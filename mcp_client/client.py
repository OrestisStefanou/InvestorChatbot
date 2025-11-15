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
        
        result = await client.call_tool(
            "getStockOverview", 
            {
                "stock_symbol": "MSFT", 
                #"order_id": "c4e372d9-04b3-4013-af44-7a56502dd40e",
            }
        )        
        print(result.structured_content)


asyncio.run(main())