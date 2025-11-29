import asyncio
from unittest.mock import MagicMock, AsyncMock
from agent_service.user_context import MongoDBUserContextService, UserPortfolioHolding

async def test_create_user_context():
    # Mock Mongo Client
    mock_client = MagicMock()
    mock_db = MagicMock()
    mock_collection = MagicMock()
    
    # Setup async insert_one
    mock_collection.insert_one = AsyncMock()
    
    mock_client.__getitem__.return_value = mock_db
    mock_db.__getitem__.return_value = mock_collection
    
    service = MongoDBUserContextService(mock_client, "test_db", "test_collection")
    
    user_id = "test_user_123"
    
    # Test case 1: Default values
    print("Testing create_user_context with defaults...")
    context = await service.create_user_context(user_id)
    assert context.userid == user_id
    assert context.userprofile == {}
    assert context.userportfolio == []
    mock_collection.insert_one.assert_called()
    print("Test case 1 passed!")

    # Test case 2: Provided values
    print("Testing create_user_context with provided values...")
    profile = {"age": 30}
    portfolio = [UserPortfolioHolding(assetclass="Equity", symbol="AAPL", name="Apple", quantity=10)]
    context = await service.create_user_context(user_id, user_profile=profile, user_portfolio=portfolio)
    assert context.userid == user_id
    assert context.userprofile == profile
    assert context.userportfolio == portfolio
    print("Test case 2 passed!")

if __name__ == "__main__":
    asyncio.run(test_create_user_context())
