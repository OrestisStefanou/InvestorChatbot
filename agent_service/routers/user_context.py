import http

from fastapi import (
    APIRouter, 
    Depends, 
    HTTPException
)
from pydantic import BaseModel
from pymongo import AsyncMongoClient

from agent_service.services.user_context import (
    MongoDBUserContextService, 
    UserContext,
    UserPortfolioHolding,
    UserContextAlreadyExistsError
)
from agent_service.config import settings
from agent_service.dependencies import get_db_client

router = APIRouter()

class UserPortfolioHoldingSchema(BaseModel):
    asset_class: str
    symbol: str
    name: str
    quantity: int


class UserContextSchema(BaseModel):
    user_id: str
    user_profile: dict | None = None
    user_portfolio: list[UserPortfolioHoldingSchema] | None = None


@router.post("/user_context", response_model=UserContextSchema, status_code=http.HTTPStatus.CREATED)
async def create_user_context(request: UserContextSchema, db_client: AsyncMongoClient = Depends(get_db_client)):
    user_context_service = MongoDBUserContextService(
        mongo_client=db_client,
        db_name=settings.MONGO_DB_NAME,
        collection_name=settings.USER_CONTEXT_COLLECTION_NAME,
    )

    # Convert UserContextSchema to UserContext
    user_context = UserContext(
        userid=request.user_id,
        userprofile=request.user_profile or {},
        userportfolio=[UserPortfolioHolding(
            assetclass=holding.asset_class,
            symbol=holding.symbol,
            name=holding.name,
            quantity=holding.quantity,
        ) for holding in request.user_portfolio or []
        ],
    )

    try:
        created_user_context = await user_context_service.create_user_context(
            user_id=request.user_id,
            user_profile=request.user_profile,
            user_portfolio=request.user_portfolio,
        )
    except UserContextAlreadyExistsError as e:
        raise HTTPException(status_code=http.HTTPStatus.CONFLICT, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=http.HTTPStatus.INTERNAL_SERVER_ERROR, detail=str(e))
    
    return UserContextSchema(
        user_id=created_user_context.userid,
        user_profile=created_user_context.userprofile,
        user_portfolio=[UserPortfolioHoldingSchema(
            asset_class=holding.assetclass,
            symbol=holding.symbol,
            name=holding.name,
            quantity=holding.quantity,
        ) for holding in created_user_context.userportfolio],
    )
        