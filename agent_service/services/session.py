from abc import ABC, abstractmethod
from enum import Enum
import uuid

from pydantic import BaseModel
from pymongo import AsyncMongoClient

class MessageRole(str, Enum):
    USER = "user"
    AGENT = "agent"


class Message(BaseModel):
    role: MessageRole
    content: str


class Session(BaseModel):
    sessionID: str
    user_id: str
    messages: list[Message]


class SessionService(ABC):
    @abstractmethod
    async def create_session(self, user_id: str) -> Session:
        pass

    @abstractmethod
    async def get_session(self, session_id: str) -> Session | None:
        pass

    @abstractmethod
    async def add_message(self, session_id: str, message: Message) -> Session | None:
        pass


class SessionNotFoundError(Exception):
    pass


class MongoDBSessionService(SessionService):
    def __init__(self, mongo_client: AsyncMongoClient, db_name: str, collection_name: str):
        self.db = mongo_client[db_name]
        self.collection = self.db[collection_name]

    async def create_session(self, user_id: str) -> Session:
        session_id = str(uuid.uuid4())
        session = Session(sessionID=session_id, user_id=user_id, messages=[])
        await self.collection.insert_one(session.model_dump())
        return session
    
    async def get_session(self, session_id: str) -> Session | None:
        doc = await self.collection.find_one({"sessionID": session_id})
        if doc:
            return Session(**doc)

        return None

    async def add_message(self, session_id: str, message: Message) -> Session | None:
        session = await self.get_session(session_id)
        if not session:
            raise SessionNotFoundError("Session not found")

        session.messages.append(message)
        await self.collection.update_one({"sessionID": session_id}, {"$set": session.model_dump()})
        return session