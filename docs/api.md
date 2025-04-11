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

