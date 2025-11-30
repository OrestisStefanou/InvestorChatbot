from abc import ABC, abstractmethod

from config import settings
from agent_service.services.session import SessionService

class ChatService(ABC):
    @abstractmethod
    async def generate_response(self, session_id: str, message: str) -> str:
        raise NotImplementedError

class AgenticChatService(ChatService):
    def __init__(self, session_service: SessionService):
        # Use the settings to create an agent here
        # Fetch the prompt from the MCP server
        self._agent = None
        self._session_service = session_service
    
    async def generate_response(self, session_id: str, message: str) -> str:
        # Get the session
        # Generate the response
        # Store the message and response in the session
        # Return the response
        return "Hello, world!"