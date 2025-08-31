# AI Investor Assistant API

## Key Benefits

* 🚀 **Plug-and-Play AI Assistant** – No need to build AI models or pipelines; simply integrate and provide immediate value.
* 🧠 **Personalized Investor Insights** – Responses can be tailored based on user context, such as portfolio composition, risk profile, and interests.
* 🤖 **Advanced AI Models** – Behind the scenes, the service leverages **OpenAI** and **Google Gemini** (configurable per use case) for high-quality answers.
* 📊 **Trusted Market Data Sources** – All responses are generated using **institutional-grade financial data** for accuracy and reliability.
* 🔍 **Continuous Quality Evaluation** – All questions and responses are securely stored to monitor performance and improve answer quality over time.
* 💡 **Increased User Engagement & Retention** – Real-time, interactive investment guidance keeps users engaged and returning.
* ⏱ **Time & Cost Savings** – Avoid the cost of developing, training, and maintaining large language models and financial knowledge bases.
* 📚 **Enhanced Support Capabilities** – AI can handle FAQs, generate follow-up questions, and guide users toward actionable insights.

---

## Example End-to-End Flow

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

## Config settings
[Runtime config settings](docs/config.md)


## Api usage examples
[API usage examples](docs/api_request_examples.md)

## Complete api reference
[API reference](docs/api.md)


## 📞 Contact
For issues, reach out via GitHub Issues or email: [stefanouorestis@gmail.com](mailto:stefanouorestis@gmail.com)

