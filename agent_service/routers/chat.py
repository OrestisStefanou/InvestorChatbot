import http

from fastapi import APIRouter, Depends
from pydantic import BaseModel
from pymongo import AsyncMongoClient

from agent_service.services.session import MongoDBSessionService
from agent_service.services.chat import AgenticChatService
from agent_service.config import settings
from agent_service.dependencies import get_db_client

router = APIRouter()


class ChatRequest(BaseModel):
    session_id: str
    message: str


class ChatResponse(BaseModel):
    response: str


@router.post("/chat", response_model=ChatResponse)
async def chat(request: ChatRequest, db_client: AsyncMongoClient = Depends(get_db_client)):
    session_service = MongoDBSessionService(
        mongo_client=db_client,
        db_name=settings.MONGO_DB_NAME,
        collection_name=settings.SESSION_COLLECTION_NAME,
    )
    chat_service = AgenticChatService(session_service)
    # TODO: CHECK FOR ANY EXCEPTIONS HERE
    response = await chat_service.generate_response(request.session_id, request.message)
    return ChatResponse(response=response)
