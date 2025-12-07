from fastapi import (
    Request,
    HTTPException,
)
from langchain_mcp_adapters.client import MultiServerMCPClient
from agent_service.config import settings

def get_db_client(request: Request):
    if not hasattr(request.app.state, "mongodb_client"):
        raise HTTPException(status_code=500, detail="Database not initialized")
    return request.app.state.mongodb_client

def get_mcp_client():
    mcp_server_client = MultiServerMCPClient(
        {
            "investing_data_tools": {
                "transport": "streamable_http",  # HTTP-based remote server
                "url": settings.MCP_SERVER_URL,
            }
        }
    )    
    return mcp_server_client
