from abc import ABC, abstractmethod

from langchain.agents import create_agent
from langchain.agents.structured_output import ToolStrategy
from langchain_mcp_adapters.client import MultiServerMCPClient
from langchain_openai import ChatOpenAI
from langchain_google_genai import ChatGoogleGenerativeAI
from langchain.messages import (
    HumanMessage,
    AIMessage, 
)
from pydantic import BaseModel

from agent_service.config import (
    settings,
    LLMProvider,
)
from agent_service.services.session import (
    Session,
    SessionService,
    SessionNotFoundError,
    MessageRole,
    Message,
)


class ResponseFormat(BaseModel):
    response: str

class ChatService(ABC):
    @abstractmethod
    async def generate_response(self, session_id: str, message: str) -> str:
        raise NotImplementedError


# TODO;
# 1. Add middleware for tool error handling
# 2. Move agent ot a separate class
#    - Pass it as a parameter to the constructor
#    - Make it a dependency on the endpoint?

class AgenticChatService(ChatService):
    def __init__(self, session_service: SessionService):
        # Use the settings to create an agent here
        # Fetch the prompt from the MCP server
        self._agent = None
        self._session_service = session_service
    
    async def generate_response(self, session_id: str, message: str) -> str:
        # We can't create the agent in the constructor because it requires async code
        if self._agent is None:
            await self._set_up_agent()
        
        # Get the session
        session = await self._session_service.get_session(session_id)
        if not session:
            raise SessionNotFoundError(f"Session {session_id} not found")
        
        # Create the messages from the session
        messages = self._create_messages_from_session(session)
        # Add the new message to conversation passed for inference
        messages.append({"role": "user", "content": message})

        # Generate the response
        print("MESSAGES:", {"messages": messages})
        response = await self._agent.ainvoke({"messages": messages})
        structured_response = response["structured_response"].response

        print("RESPONSE IS:", response)

        # Store the message and response in the session
        await self._session_service.add_message(session_id, Message(role=MessageRole.USER, content=message))
        await self._session_service.add_message(session_id, Message(role=MessageRole.AGENT, content=structured_response))
        # Return the response
        return structured_response


    def _create_messages_from_session(self, session: Session) -> list[dict[str, str]]:
        # TODO: Add some limit here(should get the limit from the settings)
        messages = []
        for message in session.messages:
            if message.role == MessageRole.USER:
                messages.append({"role": "user", "content": message.content})
            elif message.role == MessageRole.AGENT:
                messages.append({"role": "assistant", "content": message.content})
        return messages

    async def _set_up_agent(self):
        mcp_server_client = MultiServerMCPClient(
            {
                "investing_data_tools": {
                    "transport": "streamable_http",  # HTTP-based remote server
                    "url": settings.MCP_SERVER_URL,
                }
            }
        )
        tools = await mcp_server_client.get_tools()

        match settings.LLM_PROVIDER:
            case LLMProvider.OPENAI:
                model = ChatOpenAI(
                    api_key=settings.OPENAI_API_KEY,
                    model=settings.LLM_MODEL,
                    temperature=settings.TEMPERATURE,
                )
            case LLMProvider.GOOGLE:
                print(f"CREATING GOOGLE MODEL WITH MODEL: {settings.LLM_MODEL} and temperature: {settings.TEMPERATURE}")
                model = ChatGoogleGenerativeAI(
                    google_api_key=settings.GOOGLE_API_KEY,
                    model=settings.LLM_MODEL,
                    temperature=settings.TEMPERATURE,
                )
            case _:
                raise ValueError(f"Unknown LLM provider: {settings.LLM_PROVIDER}")
        
        self._agent = create_agent(
            model, 
            tools=tools,
            response_format=ToolStrategy(ResponseFormat),
            system_prompt="You are a helpful investment assistant. Use the tools provided to answer the user's question.",   # TODO: Fetch the prompt from the MCP server
        )
