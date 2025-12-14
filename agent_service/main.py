from contextlib import asynccontextmanager

from fastapi import FastAPI, Depends, HTTPException
from pymongo import AsyncMongoClient

from agent_service.routers import (
    session,
    user_context,
    chat,
)
from agent_service.config import settings

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    app.state.mongodb_client = AsyncMongoClient(settings.MONGO_URI)
    yield
    # Shutdown
    await app.state.mongodb_client.close()


app = FastAPI(lifespan=lifespan)

app.include_router(session.router)
app.include_router(user_context.router)
app.include_router(chat.router)