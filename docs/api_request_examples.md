# 📘 API Usage Examples

---

## 📰 How to ask a questions about investing education

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "What are the benefits of investing?",
  "topic": "education",
  "session_id": "<session_id>",
  "topic_tags": {}
}'
```


## 📰 How to ask a questions about all the stock sectors

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Which are the best performing sectors?",
  "topic": "sectors",
  "session_id": "<session_id>",
  "topic_tags": {}
}'
```

## 📰 How to ask a question about a specific sector

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Which are the top 5 stocks in this sector?",
  "topic": "sectors",
  "session_id": "<session_id>",
  "topic_tags": {
    "sector_name": "technology"
  }
}'
```

## 📰 How to ask a general questions about a specific stock (Microsoft in this example(symbol=MSFT))

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Can you give me an overview of the stock?",
  "topic": "stock_overview",
  "session_id": "<session_id>",
  "topic_tags": {
    "stock_symbols": ["MSFT"]
  }
}'
```

## 📰 How to ask about a high level comparison of two stocks(Microsoft and Meta in this example) 

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Can you perform a high level comparison between the two stocks?",
  "topic": "stock_overview",
  "session_id": "<session_id>",
  "topic_tags": {
    "stock_symbols": ["MSFT", "META"]
  }
}'
```

## 📰 How to ask a question about the balance sheet of a stock(Microsoft this example) 

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Can you analyze the balance sheet?",
  "topic": "stock_financials",
  "session_id": "<session_id>",
  "topic_tags": {
    "stock_symbols": ["MSFT"],
    "balance_sheet": true
  }
}'
```

## 📰 How to ask a question to compare the balance sheet of two stocks(Microsoft and Meta in this example) 

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Can you compare the balance sheets of the two stocks?",
  "topic": "stock_financials",
  "session_id": "<session_id>",
  "topic_tags": {
    "stock_symbols": ["MSFT", "META"],
    "balance_sheet": true
  }
}'
```

## 📰 How to ask a question about the income statement of a stock (Microsoft in this example)

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Can you analyze the income statement?",
  "topic": "stock_financials",
  "session_id": "<session_id>",
  "topic_tags": {
    "stock_symbols": ["MSFT"],
    "income_statement": true
  }
}'
```

## 📰 How to ask a question to compare the income statements of two stocks (Microsoft and Meta in this example)

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Can you compare the income statements of the two stocks?",
  "topic": "stock_financials",
  "session_id": "<session_id>",
  "topic_tags": {
    "stock_symbols": ["MSFT", "META"],
    "income_statement": true
  }
}'
```

## 📰 How to ask a question about the cash flow statement of a stock (Microsoft in this example)

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Can you analyze the cash flow statement?",
  "topic": "stock_financials",
  "session_id": "<session_id>",
  "topic_tags": {
    "stock_symbols": ["MSFT"],
    "cash_flow": true
  }
}'
```

## 📰 How to ask a question to compare the cash flow statements of two stocks (Microsoft and Meta in this example)

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Can you compare the cash flow statements of the two stocks?",
  "topic": "stock_financials",
  "session_id": "<session_id>",
  "topic_tags": {
    "stock_symbols": ["MSFT", "META"],
    "cash_flow": true
  }
}'
```

## 📰 How to ask a question about ETFs

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "What are the key advantages of investing in ETFs?",
  "topic": "etfs",
  "session_id": "<session_id>",
  "topic_tags": {}
}'
```

## 📰 How to ask a question about a specific ETF (Goldman Sachs Physical Gold ETF in this example with etf_symbol=AAAU)

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Which are the holdings of this ETF?",
  "topic": "etfs",
  "session_id": "<session_id>",
  "topic_tags": {
    "etf_symbols": ["AAAU"]
  }
}'
```


## 📰 How to get a summary of the latest market news for NVDA

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Can you summarize the market news?",
  "topic": "news",
  "session_id": "<session_id>",
  "topic_tags": {
    "stock_symbols": ["NVDA"]
  }
}'
```

## 📰 How to get a summary of the latest market news

### 🔧 Request

```bash
curl --location 'http://localhost:1323/chat' \
--header 'Content-Type: application/json' \
--data '{
  "question": "Can you summarize the market news?",
  "topic": "news",
  "session_id": "<session_id>",
  "topic_tags": {}
}'
```

---
