from contextlib import asynccontextmanager

from fastapi import FastAPI, Depends, HTTPException
from pymongo import AsyncMongoClient

from agent_service.routers import session
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