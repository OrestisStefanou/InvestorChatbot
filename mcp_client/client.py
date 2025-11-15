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
            "getStockFinancials", 
            {
                "stock_symbol": "APM", 
                "include_balance_sheets": True,
                "include_income_statements": True,
                "include_cash_flows": True,
                "limit": 1,
            }
        )        
        print(result.structured_content)


asyncio.run(main())