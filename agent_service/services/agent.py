from langchain_mcp_adapters.client import MultiServerMCPClient
from langchain.agents import create_agent
from langchain_openai import ChatOpenAI
from langchain_google_genai import ChatGoogleGenerativeAI
from langchain.agents.structured_output import ToolStrategy
from langchain.tools import BaseTool
from langchain.chat_models import BaseChatModel
from langchain.agents.middleware import (
    AgentMiddleware,
    wrap_tool_call,
)
from langchain.messages import ToolMessage
from langchain_anthropic import ChatAnthropic
from pydantic import BaseModel

from agent_service.services.session import (
    MessageRole,
    Message,
)
from agent_service.config import (
    settings,
    LLMProvider,
)

class ToolErrorMiddleware(AgentMiddleware):
    async def awrap_tool_call(self, request, handler):
        try:
            return await handler(request)
        except Exception as e:
            return ToolMessage(
                content=f"Tool error: ({str(e)})",
                tool_call_id=request.tool_call["id"],
            )


class TextResponseFormat(BaseModel):
    response: str


class Agent:
    def __init__(
        self, 
        tools: list[BaseTool],
        model: BaseChatModel,
        response_format: ToolStrategy,
        system_prompt: str,
        middleware: list[AgentMiddleware],
    ):
        self._agent= create_agent(
            model, 
            tools=tools,
            response_format=response_format,
            system_prompt=system_prompt,
            middleware=middleware,
        )
    
    async def generate_response(self, conversation: list[Message]) -> str:
        messages = []
        # Keep the last settings.CONVERSATION_MESSAGES_LIMIT messages
        if len(conversation) > settings.CONVERSATION_MESSAGES_LIMIT:
            conversation = conversation[-settings.CONVERSATION_MESSAGES_LIMIT:] 
        
        for message in conversation:
            if message.role == MessageRole.USER:
                messages.append({"role": "user", "content": message.content})
            elif message.role == MessageRole.AGENT:
                messages.append({"role": "assistant", "content": message.content})
        
        response = await self._agent.ainvoke({"messages": messages})
        return response["structured_response"].response


# TODO: WE SHOULD GET THIS FROM THE MCP SERVER
prompt = """
You are a professional investment advisor of a client with user_id = {user_id}. Your job is to answer to any investing related questions and ask anything that you think would be useful to know  about your client to give the best personalised investing advice. 
ALWAYS follow the instructions below:
# INSTRUCTIONS
1. Always use getUserContext tool to get your user's context in order to make your responses as personalised  as possible (Do this in the background, don't let the user know that you are fetching their information to make it look like you already know it)
2. Use the updateUserContext tool to store any information about the user(your client) that you think will be useful to have for the future(don't ask the user for permission to do this, think about this as your personal notes about the user to help you give more personalised answers).
3. You should try to obtain the following information about the user(and anything else that you think would be useful):
    - The user's age
    - The user's investing knowledge level (beginner, intermediate, advanced)
    - The user's investment goals
    - The user's risk tolerance
    - The user's investment time horizon
    - The user's current investment portfolio
4. Your should use your existing tools to provide your answers if possible.
5. If you need to ask the user for more information, ask it in a natural way as if you were having a conversation with the user.
6. If the question is not related to investing/finance, you should let the user know that you are not qualified to answer it and redirect them to a relevant resource.
"""


class AgentService:
    def __init__(
        self, 
        mcp_client: MultiServerMCPClient, 
    ):
        self._mcp_client = mcp_client

    async def generate_response(
        self,
        user_id: str, 
        conversation: list[Message],
        response_format: BaseModel,
    ) -> str:
        agent = await self._create_agent(prompt.format(user_id=user_id), response_format)
        response = await agent.generate_response(conversation)
        return response

    async def _create_agent(self, system_prompt: str, response_format: BaseModel) -> Agent:
        tools = await self._mcp_client.get_tools()
        match settings.LLM_PROVIDER:
            case LLMProvider.OPENAI:
                model = ChatOpenAI(
                    api_key=settings.OPENAI_API_KEY,
                    model=settings.LLM_MODEL,
                    temperature=settings.TEMPERATURE,
                )
            case LLMProvider.GOOGLE:
                model = ChatGoogleGenerativeAI(
                    google_api_key=settings.GOOGLE_API_KEY,
                    model=settings.LLM_MODEL,
                    temperature=settings.TEMPERATURE,
                )
            case LLMProvider.ANTHROPIC:
                model = ChatAnthropic(
                    api_key=settings.ANTHROPIC_API_KEY,
                    model=settings.LLM_MODEL,
                    temperature=settings.TEMPERATURE,
                )
            case _:
                raise ValueError(f"Unknown LLM provider: {settings.LLM_PROVIDER}") 

        return Agent(
            tools=tools,
            model=model,
            response_format=ToolStrategy(response_format),
            system_prompt=system_prompt,
            middleware=[ToolErrorMiddleware()],
        )
