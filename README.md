# 💹 AI Investment Copilot API

An open-source **AI-powered investment assistant API** designed to provide intelligent, context-aware financial insights, portfolio analysis, and conversational interactions for investors and analysts.

This project offers a modular backend for an **AI investment copilot** — integrating topic extraction, financial context tagging, and user personalization — built for easy use in chatbots, dashboards, or financial research tools.

---

## 🚀 Features

* 🧠 **Conversational AI for finance** — powered by OpenAI, Gemini, or Ollama models.
* 💬 **Contextual chat sessions** — persistent session tracking for ongoing conversations.
* 📊 **Topic & tag extraction** — automatically identify topics (e.g., "stock_overview") and extract context like tickers or financial statements.
* 👤 **User personalization** — customize responses using user profiles and portfolios.
* 🔎 **Dynamic FAQ & sector data** — retrieve FAQs, tickers, sectors, and ETFs for market insights.
* 🤖 **Follow-up question generation** — intelligently guide users toward deeper exploration.
* ⚙️ **Configurable and extensible** — easily switch between LLM or database providers using environment variables.

---

## 🏗️ Architecture Overview

The API consists of multiple endpoints handling:

* **Session management**
* **Chat completions** (streamed responses)
* **Topic and tag extraction**
* **User context**
* **Follow-up question generation**
* **FAQ and market data retrieval**

Behind the scenes, topic and tag extraction is performed in **two LLM steps**:

1. Detect the topic of the question.
2. Extract relevant tags (e.g., stock symbols, sectors, etc.) based on that topic.

> See [`topic_tag_extractor.md`](topic_tag_extractor.md) for implementation details.

---

## ⚙️ Configuration

Configuration is handled via environment variables or a `.env` file.
Key providers and settings include:

### Providers

| Type            | Options                       |
| --------------- | ----------------------------- |
| LLM Provider    | `OPEN_AI`, `OLLAMA`, `GEMINI` |
| Database        | `MONGO_DB`, `BADGER`          |
| Session Storage | `MONGO_DB`, `MEMORY`          |

### Example `.env`

```env
LlmProvider=OPEN_AI
OPEN_AI_API_KEY=your_openai_key_here
DatabaseProvider=MONGO_DB
MONGO_DB_URI=mongodb://localhost:27017
MONGO_DB_NAME=investment_copilot
FAQ_LIMIT=5
CONV_MSG_LIMIT=10
FOLLOW_UP_QUESTIONS_NUM=5
CACHE_TTL=3600
```

> For a full list of configuration variables, see [`config.md`](docs/config.md).

---

## 🧩 API Overview

### 🔹 **Session Management**

* `POST /session` – Create a new session.
* `GET /session/:session_id` – Retrieve conversation history.

### 🔹 **Chat & AI Responses**

* `POST /chat` – Generate streamed chat responses.
* `POST /chat/extract_topic_and_tags` – Extract the topic and financial tags from a question.

### 🔹 **User Context**

* `POST /user_context` – Create a personalized user profile and portfolio.
* `PUT /user_context` – Update user context.
* `GET /user_context/:user_id` – Retrieve existing user context.

### 🔹 **Follow-Up Questions**

* `POST /follow_up_questions` – Generate next-step questions to continue engagement.

### 🔹 **Market Data & FAQs**

* `GET /faq` – Retrieve FAQs for a specific topic.
* `GET /topics` – List all supported FAQ topics.
* `GET /tickers` – Search and list stock tickers.
* `GET /sectors` – Retrieve sector-level data.
* `GET /sectors/stocks/:sector` – Get all stocks in a specific sector.
* `GET /etfs` – Retrieve a list of ETFs.

> Detailed request and response formats are available in [`api.md`](docs/api.md).

---

## 🧠 Topic & Tag Extraction

Topic and tag extraction is handled in **two steps**:

1. **Topic Extraction** — Identify the main subject of the user’s query.
2. **Tag Extraction** — Extract related financial entities (e.g., `stock_symbols`, `sector_name`, etc.) based on topic-specific prompts.

### Advantages

* Modular, reusable prompts for each topic.
* Improved LLM accuracy and focus per task.

### Future Improvements

* Use embeddings and vector similarity to limit symbol search space for improved efficiency.

For more, see [`topic_tag_extractor.md`](dosc/topic_tag_extractor.md).

---

## 🧰 Installation & Setup

### Prerequisites

* Go 1.21+
* MongoDB (if using `MONGO_DB`)
* Python (for optional Streamlit client)

### Install Dependencies

```bash
go mod tidy
```

### Run the API Server

```bash
make run_investbot
```

Server will start at:

```
http://localhost:1323
```

---

## 💻 Example Client (Optional)

A minimal Streamlit client is available for testing:

```bash
pip install streamlit requests
streamlit run client.py
```

---

## 🧪 Example Use Flow

Below is a quick way to test the API using simple HTTP requests. You can use **curl**, **Postman**, or any HTTP client.

### Step 1 – Create a Session

```sh
POST /session
```

**Response:**

```json
{ "session_id": "abc123xyz" }
```

---

### Step 2 – Create a User (Optional)

```sh
POST /user_context
Content-Type: application/json

{
  "user_id": "user_123",
  "user_profile": {
    "risk_tolerance": "moderate",
    "investment_horizon": "5 years"
  },
  "user_portfolio": [
    { "asset_class": "stock", "symbol": "AAPL", "quantity": 15, "portfolio_percentage": 50 },
    { "asset_class": "etf", "symbol": "SPY", "quantity": 5, "portfolio_percentage": 50 }
  ]
}
```

---

### Step 3 – Extract Topic & Tags (Optional)

```sh
POST /chat/extract_topic_and_tags
Content-Type: application/json

{
  "question": "How did Apple perform last quarter?",
  "session_id": "abc123xyz"
}
```

**Example Response:**

```json
{
  "topic": "stock_overview",
  "topic_tags": {
    "stock_symbols": ["AAPL"]
  }
}
```

---

### Step 4 – Generate a Chat Response

```sh
POST /chat
Content-Type: application/json

{
  "question": "How did Apple perform last quarter?",
  "topic": "stock_overview",
  "session_id": "abc123xyz",
  "topic_tags": {
    "stock_symbols": ["AAPL"]
  }
}
```

Response will stream back the AI-generated answer.

---

### Step 5 – Generate Follow-Up Questions

```sh
POST /follow_up_questions
Content-Type: application/json

{
  "session_id": "abc123xyz",
  "number_of_questions": 3
}
```

**Example Response:**

```json
{
  "follow_up_questions": [
    "Would you like to see a breakdown of Apple's revenue?",
    "Should I compare Apple’s results to Microsoft?",
    "Do you want to explore Apple’s future growth potential?"
  ]
}
```
---

## How to test it with simple Streamlit ui app

### ⚙️ Installation
**Install dependencies**

   ```bash
   pip install streamlit requests
   ```

---

### ▶️ Run the App

1. Make sure your backend API is running at:

   ```
   http://localhost:1323
   ```

2. Start the Streamlit app:

   ```bash
   streamlit run client.py
   ```

3. Open the link shown in your terminal (usually `http://localhost:8501`) to use the chatbot.
4. For personalised responses edit the [`client.py`](ui/client.py) file and update the USER_ID global variable with the id of your user.
---

## Next steps
1. Add support for crypto
2. Create an mcp server to expose the tools that the llms are using in the background.

## 🛡️ License

This project is open source and distributed under the **MIT License**.
Feel free to fork, contribute, and build your own investment copilots.

---

## 🤝 Contributing

Contributions are welcome!
If you’d like to improve prompts, extend APIs, or optimize topic extraction:

1. Fork the repo
2. Create a feature branch
3. Submit a pull request

---

## 📞 Contact
For issues, reach out via GitHub Issues or email: [stefanouorestis@gmail.com](mailto:stefanouorestis@gmail.com)

