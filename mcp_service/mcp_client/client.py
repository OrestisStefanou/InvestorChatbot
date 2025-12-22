# THIS IS AN MCP CLIENT JUST TO TEST THE TOOLS RESPONSES

import uuid

import asyncio
from fastmcp import Client, FastMCP
from datetime import datetime, timedelta


# HTTP server
client = Client("http://127.0.0.1:8080/mcp")

async def main():
    async with client:
        # Basic server interaction
        await client.ping()

        # Current datetime
        now = datetime.now()

        # Datetime 5 days ago
        five_days_ago = now - timedelta(days=5)
        
        # result = await client.call_tool(
        #     name="getCommodityTimeSeries", 
        #     arguments={
        #         'commodity_name': 'CrudeOil', 
        #         'limit': 10,
        #         # 'user_profile': {'age': 28, 'investment_knowledge_level': 'Intermediate', 'name': 'Orestis', 'risk_apettite': 'medium'}, 
        #         # 'user_portfolio': [
        #         #     {'asset_class': 'stock', 'symbol': '', 'name': 'META', 'quantity': 20, 'portfolio_percentage': 0.25}, 
        #         #     {'asset_class': 'stock', 'symbol': '', 'name': 'Microsoft', 'quantity': 15, 'portfolio_percentage': 0.25},
        #         #     {'asset_class': 'crypto', 'symbol': 'BTC', 'name': 'Bitcoin', 'quantity': 0.01, 'portfolio_percentage': 0.5}
        #         # ]
        #     }
        # )
        
        result = await client.get_prompt(
            name="investment_advisor_prompt", 
            arguments={
                'user_id': str(uuid.uuid4()),
            }
        )        
        print(result)


asyncio.run(main())