from abc import ABC, abstractmethod

from pydantic import BaseModel
from pymongo import AsyncMongoClient

# The naming of the fields is done to match mcp server schema
class UserPortfolioHolding(BaseModel):
    assetclass: str
    symbol: str
    name: str
    quantity: int

class UserContext(BaseModel):
    userid: str
    userprofile: dict
    userportfolio: list[UserPortfolioHolding]


class UserContextAlreadyExistsError(Exception):
    pass


class UserContextService(ABC):
    @abstractmethod
    async def create_user_context(
        self, 
        user_id: str,
        user_profile: dict | None = None,
        user_portfolio: list[UserPortfolioHolding] | None = None,
    ) -> UserContext | None:
        pass


class MongoDBUserContextService(UserContextService):
    def __init__(self, mongo_client: AsyncMongoClient, db_name: str, collection_name: str):
        self.db = mongo_client[db_name]
        self.collection = self.db[collection_name]

    async def create_user_context(
        self, 
        user_id: str,
        user_profile: dict | None = None,
        user_portfolio: list[UserPortfolioHolding] | None = None,
    ) -> UserContext | None:

        # Check if user context for user_id already exists
        existing_user_context = await self.collection.find_one({"userid": user_id})
        if existing_user_context:
            raise UserContextAlreadyExistsError(f"User context already exists for user_id: {user_id}")

        user_context = UserContext(
            userid=user_id,
            userprofile=user_profile if user_profile is not None else {},
            userportfolio=user_portfolio if user_portfolio is not None else []
        )
        await self.collection.insert_one(user_context.model_dump())
        return user_context
