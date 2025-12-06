import asyncio

from langchain_mcp_adapters.client import MultiServerMCPClient
from langchain.agents import create_agent
from langchain_openai import ChatOpenAI
from langchain_google_genai import ChatGoogleGenerativeAI
from langchain.agents.structured_output import ToolStrategy
from pydantic import BaseModel


class ResponseFormat(BaseModel):
    response: str

async def main():
    client = MultiServerMCPClient(  
        {
            "investing_analysis_tools": {
                "transport": "streamable_http",  # HTTP-based remote server
                # Ensure you start your weather server on port 8000
                "url": "http://127.0.0.1:8080/mcp",
            }
        }
    )

    tools = await client.get_tools()  

    # model = ChatOpenAI(
    #     model="gpt-5",
    #     temperature=0.1,
    # )

    model = ChatGoogleGenerativeAI(
        model="gemini-2.5-flash"
    )

    agent = create_agent(
        model, 
        tools=tools,
        response_format=ToolStrategy(ResponseFormat),
    )

    result = await agent.ainvoke(
        {"messages": [{"role": "user", "content": "what are the latest market news?"}]}
    )
    print(result["structured_response"].response)

asyncio.run(main())


