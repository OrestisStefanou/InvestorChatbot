import asyncio

from pymongo import AsyncMongoClient

from session import MongoDBSessionService, Message
from user_context import MongoDBUserContextService


from config import settings

async def main():
    mongo_client = AsyncMongoClient(settings.MONGO_URI)
    user_context_service = MongoDBUserContextService(mongo_client, settings.MONGO_DB_NAME, settings.USER_CONTEXT_COLLECTION_NAME)
    user_context = await user_context_service.create_user_context("lando-norris")
    print(user_context)


asyncio.run(main())
