# Create New Session API

## Endpoint

### POST `/session`

Creates a new session and returns the session ID.

## Request Parameters

_No parameters are required in the request body or query._

## Response

### Success Response (201 Created)

#### Example Response Body:
```json
{
  "session_id": "abc123xyz"
}
```

### Error Response (500 Internal Server Error)

#### Example Response Body:
```json
{
  "error": "failed to create session"
}
```

## Notes
- This endpoint is used to create a new session on the server.
- A successful request returns a `session_id`, which can be used in the chat endpoint.
- If an error occurs during session creation, a relevant error message will be returned in the response body.

## Example Request
```sh
POST /session
```

This request would create a new session and return the newly generated session ID.

--- 

Absolutely! Here's the full API documentation for the `/chat` POST endpoint, based on the code you shared:

---

# Chat Completion API

## Endpoint

### POST `/chat`

Generates a streaming chat response based on the user's question, topic, session, and contextual tags (like stock or financial statement info).

## Request Body

| Field             | Type      | Required | Description                                                                 |
|------------------|-----------|----------|-----------------------------------------------------------------------------|
| `question`        | string    | Yes      | The user's question to be answered.                                         |
| `topic`           | string    | Yes      | The context/topic for the chat (e.g., "finance", "markets").                |
| `session_id`      | string    | Yes      | A valid session ID created via the `/session` endpoint.                     |
| `topic_tags`      | object    | No       | Optional tags to add financial context. See `Topic Tags` below.             |

### Topic Tags (Optional `topic_tags` object)

| Field              | Type    | Required | Description                                                              |
|-------------------|---------|----------|--------------------------------------------------------------------------|
| `sector_name`      | string  | No       | Name of the relevant sector (e.g., "Technology").                         |
| `industry_name`    | string  | No       | Name of the industry (e.g., "Semiconductors").                            |
| `stock_symbol`     | string  | No       | Specific stock symbol (e.g., "AAPL").                                     |
| `balance_sheet`    | boolean | No       | Whether to include balance sheet context.                                |
| `income_statement` | boolean | No       | Whether to include income statement context.                             |
| `cash_flow`        | boolean | No       | Whether to include cash flow context.                                    |
| `etf_symbol`       | string  | No       | ETF symbol if the question is related to an ETF.                          |

### Example Request Body
```json
{
  "question": "How did Apple perform last quarter?",
  "topic": "stock_overview",
  "session_id": "abc123xyz",
  "topic_tags": {
    "stock_symbol": "AAPL"
  }
}
```

## Response

### Success Response (200 OK – Streamed)

The response is a stream of JSON-encoded text chunks representing the chat reply. Each chunk is a string:
```json
"Apple reported strong earnings with increased revenue in Q4..."
```

> Note: This is streamed using server-sent events (chunked HTTP), not returned as a complete JSON object.

### Error Responses

#### 400 Bad Request
Occurs when the request payload is invalid or missing required fields.

```json
{
  "error": "question field is required"
}
```

```json
{
  "error": "session not found"
}
```

```json
{
  "error": "invalid topic"
}
```

#### 500 Internal Server Error

```json
{
  "error": "an unexpected error occurred"
}
```

## Notes
- Fields `question`, `topic`, and `session_id` are required.
- If `session_id` is invalid or expired, a 400 error will be returned.
- This endpoint returns a **streaming** response, suitable for chat UIs that render text incrementally.
- The `topic_tags` object allows for fine-grained control over the context of the AI's response, especially when discussing financials.

## Example Request
```sh
POST /chat
Content-Type: application/json

{
  "question": "What is the sentiment of the latest market news?",
  "topic": "news",
  "session_id": "abc123xyz",
  "topic_tags": {}
}
```

This request would trigger a streamed AI response about the semiconductor industry.

---


# Get Sector Stocks API

## Endpoint

### GET `/sectors/stocks/:sector`

Retrieves a list of stocks belonging to a specific sector.

## Request Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `sector`  | string | Yes      | The sector identifier used to filter stocks. This should be the `url_name` field from the `/sectors` endpoint response. |

## Response

### Success Response (200 OK)

#### Example Response Body:
```json
{
  "SectorStocks": [
    {
      "Symbol": "AAPL",
      "CompanyName": "Apple Inc.",
      "MarketCap": 2500000000000
    },
    {
      "Symbol": "MSFT",
      "CompanyName": "Microsoft Corporation",
      "MarketCap": 2200000000000
    }
  ]
}
```

### Error Response (500 Internal Server Error)

#### Example Response Body:
```json
{
  "error": "An error occurred while fetching sector stocks"
}
```

## Notes
- The `sector` parameter should be a valid `url_name` from the `/sectors` endpoint response (e.g., `technology`, `finance`).
- The response returns an array of stock objects, each containing:
  - `Symbol`: The stock ticker symbol.
  - `CompanyName`: The name of the company.
  - `MarketCap`: The company's market capitalization.
- If an error occurs while fetching the stocks, an appropriate error message will be returned.

## Example Request
```sh
GET /sectors/stocks/technology
```

This request would return a list of technology sector stocks.

---

# Get Sectors API

## Endpoint

### GET `/sectors`

Retrieves a list of all available sectors and their details.

## Response

### Success Response (200 OK)

#### Example Response Body:
```json
{
  "Sectors": [
    {
      "name": "Technology",
      "url_name": "technology",
      "number_of_stocks": 150,
      "market_cap": 15000000000000,
      "dividend_yield_pct": 1.5,
      "pe_ratio": 25.4,
      "profit_margin_pct": 12.3,
      "one_year_change_pct": 15.2
    },
    {
      "name": "Finance",
      "url_name": "finance",
      "number_of_stocks": 120,
      "market_cap": 12000000000000,
      "dividend_yield_pct": 2.1,
      "pe_ratio": 18.7,
      "profit_margin_pct": 10.5,
      "one_year_change_pct": 8.4
    }
  ]
}
```

### Error Response (500 Internal Server Error)

#### Example Response Body:
```json
{
  "error": "An error occurred while fetching sectors"
}
```

## Notes
- The response returns an array of sector objects, each containing:
  - `name`: The name of the sector.
  - `url_name`: A URL-friendly version of the sector name.
  - `number_of_stocks`: The total number of stocks in the sector.
  - `market_cap`: The total market capitalization of the sector.
  - `dividend_yield_pct`: The average dividend yield percentage.
  - `pe_ratio`: The average price-to-earnings ratio.
  - `profit_margin_pct`: The average profit margin percentage.
  - `one_year_change_pct`: The percentage change in sector value over the past year.

## Example Request
```sh
GET /sectors
```

This request would return a list of all available sectors and their details.

