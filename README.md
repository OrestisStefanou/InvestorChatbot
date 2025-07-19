# 💼 AI Investor Assistant API – Powering the Next Generation of Financial Intelligence

## 🌟 Overview

**The AI Investor Assistant API is a plug-and-play SaaS backend** that empowers any financial product—be it an investment app, trading platform, portfolio tracker, or finance education tool—to offer intelligent, **personalized investor guidance** via natural language chat.

Built to seamlessly integrate into existing systems, this API delivers deep financial insights by combining:
- Real-time user portfolio data,
- Sector & stock-level context,
- State-of-the-art LLM technology.

Whether you're looking to **engage users**, **reduce churn**, **increase retention**, or **stand out with AI-driven features**, this toolkit gives your app a serious competitive edge—with **zero AI infrastructure overhead.**

---

## 🚀 Why It Matters for Your App

### 🧠 Transform User Experience
Engage your users with smart, context-aware financial conversations—just like having a personal financial assistant inside your app.

### 📊 Personalization That Scales
Tailor every answer based on each user’s risk tolerance, age, portfolio composition, and investing goals. Go beyond generic responses.

### 🔁 Increase Engagement & Retention
Use automated follow-up questions, personalized investment education, and FAQ discovery to keep users exploring and coming back.

### 💸 Save Time & Resources
Skip the costly R&D of building your own AI stack. With this API, you get a prebuilt, production-ready AI assistant tailored for finance.

---

## 💡 What Can It Do?

Here’s what your app unlocks with the AI Investor Assistant:

| Capability                    | Description                                                                 |
|------------------------------|-----------------------------------------------------------------------------|
| 🗣️ Natural Language Q&A       | Users can ask things like "How is Tesla doing?" or "What are ETFs?"         |
| 📈 Portfolio-Based Advice     | Responses factor in the user's own holdings and investment profile          |
| 🧠 Topic + Tag Detection      | Automatically understands context: sectors, tickers, statements, etc.       |
| 📚 Smart Financial Education  | Returns curated FAQs and bite-sized financial lessons                       |
| 🔄 Follow-Up Questions        | Generates intelligent, contextual follow-ups to keep the conversation going |
| ⚡ Real-Time Streaming        | Chat replies stream in real time, ideal for conversational UI               |
| 🔌 Easy Integration           | Just a few API calls—no need to host models or fine-tune prompts            |

---

## 🧱 How It Works – Under the Hood

The system is composed of modular, interoperable endpoints that handle everything from session management to chat to follow-up generation.

### 1. Start a Session
Every conversation starts with:
```http
POST /session
````

Returns a `session_id` used across endpoints to track chat history.

---

### 2. Optional: Provide User Context

Add user-specific data like portfolio and profile (age, risk level, etc.):

```http
POST /user_context
```

Now all replies will adapt to this context. The assistant knows what the user holds and who they are.

---

### 3. Extract the Topic & Financial Tags

Not sure how to interpret a user query like:

> “How did Microsoft and Apple perform in the tech sector?”

Let the API figure it out:

```http
POST /chat/extract_topic_and_tags
```

Returns:

```json
{
  "topic": "stock_overview",
  "topic_tags": {
    "stock_symbols": ["AAPL", "MSFT"],
    "sector_name": "Technology"
  }
}
```

---

### 4. Generate the Response

Send the enriched request to the main chat endpoint:

```http
POST /chat
```

The assistant will stream back a detailed, human-like response incorporating everything it knows—context, holdings, market data, and tone.

---

### 5. Drive Engagement With Follow-Ups

Want to keep the user engaged and exploring?

```http
POST /follow_up_questions
```

Returns a list of smart, personalized follow-up questions like:

* “Would you like to compare Microsoft’s and Apple’s R\&D spend?”
* “Should I include ETF alternatives in this analysis?”

---

### 6. Provide Learning Resources

Return beginner-friendly educational content by topic:

```http
GET /faq?faq_topic=income_statement
```

Great for onboarding new investors and reinforcing trust in your platform.

---

## 🧪 Developer-First Example Workflow

```bash
# Create a session
POST /session

# Send a user query
POST /chat/extract_topic_and_tags
{
  "question": "Tell me about Apple and Microsoft’s financials.",
  "session_id": "abc123xyz"
}

# Use topic + tags to fetch the response
POST /chat
{
  "question": "Tell me about Apple and Microsoft’s financials.",
  "topic": "stock_overview",
  "session_id": "abc123xyz",
  "topic_tags": {
    "stock_symbols": ["AAPL", "MSFT"]
  }
}
```

---

## 🛠 Configuration and Setup

Supports both **OpenAI** and **Ollama** as model providers. Just set the `.env` like so:

```env
LLM_PROVIDER=OPEN_AI
OPEN_AI_API_KEY=your-key
OPEN_AI_MODEL_NAME=gpt-4o-mini

# or for local LLMs via Ollama
LLM_PROVIDER=OLLAMA
OLLAMA_BASE_URL=http://localhost:11434
OLLAMA_MODEL_NAME=llama3.2
```

Other tunable settings include:

* `FAQ_LIMIT`: Number of FAQ items per topic
* `CONV_MSG_LIMIT`: How many past messages to retain in chat memory
* `BASE_LLM_TEMPERATURE`: Controls creativity vs. precision

---

## 🔍 How Topic/Tag Extraction Works

The system performs this in **two stages**:

1. **Topic Extraction**: LLM analyzes the user’s question to identify the domain (e.g. "stock\_overview").
2. **Tag Extraction**: Based on the topic, a separate prompt extracts relevant symbols, sectors, or financial docs.

✅ Modular and accurate
⚠️ Slightly more expensive due to two LLM calls (unless topic requires no tags)

See [topic\_tag\_extractor.md](./topic_tag_extractor.md) for deep details.

---

## 🌍 Real-World Use Cases

### 💸 Investment Platforms

Add a built-in assistant that answers portfolio questions, explains stock movement, and suggests strategic rebalancing.

### 📊 Portfolio Apps

Offer daily summaries, risk analysis, or chat-based insights on held positions.

### 🏫 Financial Literacy Tools

Deliver educational Q\&A, contextual FAQs, and guided learning for beginners.

### 📈 Market News Aggregators

Transform headlines into plain-English explanations with sentiment and stock impact breakdowns.

---

## 🧩 Available API Endpoints

| Feature                   | Endpoint                            |
| ------------------------- | ----------------------------------- |
| Start Session             | `POST /session`                     |
| Chat with Assistant       | `POST /chat`                        |
| Extract Topic/Tags        | `POST /chat/extract_topic_and_tags` |
| Set User Context          | `POST /user_context`                |
| Update User Context       | `PUT /user_context`                 |
| Get Follow-Up Suggestions | `POST /follow_up_questions`         |
| Retrieve FAQs             | `GET /faq?faq_topic=...`            |
| Fetch Tickers/Sectors     | `GET /tickers`, `GET /sectors`      |
| Discover ETFs             | `GET /etfs`                         |

---

## 📈 Build Smarter Finance Experiences – Today

This API isn’t just a chatbot—it's a financial brain for your application.

Instead of answering generically, it knows:

* Who the user is
* What they hold
* What they want to know
* What’s going on in the market

And it responds accordingly.

---

## 👩‍💻 Who Is It For?

* 🏦 Fintech Startups
* 📱 App Developers
* 💰 Wealth Management Platforms
* 📰 Market Analysis Tools
* 📚 Financial Education Portals

If your product has users that care about money, this assistant can help you build trust, engagement, and retention—while future-proofing your UX with AI.

---

## ✅ Get Started

1. Set your `.env` and pick your model provider.
2. Integrate `POST /session`, `POST /chat`, and optionally `/user_context`.
3. Embed follow-ups, FAQs, and personalization as needed.
4. Launch your AI investor assistant.

---

> Want to add real value to your app with smart, personalized finance insights? **Start integrating today.**

## 🛠 Installation

1. Clone the repository:
2. Navigate into the project directory:
   ```bash
   cd InvestorChatbot
   ```
3. Install dependencies:
   ```bash
   make install
   ```
4. Set up API keys in a `.env` file:
   ```bash
   OPEN_AI_API_KEY=your_api_key_here
   ```
5. Run the chatbot:
   ```bash
   make run_investbot
   ```

[Runtime config settings](docs/config.md)

## 🎯 Usage
- Start the chatbot locally:
  ```bash
  make run_investbot
  ```
- Generate a session:
  ```bash
  curl --location --request POST 'http://localhost:1323/session'
  ```
- Get a response by calling the following endpoint (streaming chunk by chunk):
  ```bash
  curl --location 'http://localhost:1323/chat' \
  --header 'Content-Type: application/json' \
  --data '{
      "question": "What are the latest market news?",
      "topic": "news",
      "session_id": "<session_id from response above>",
      "topic_tags": {}
  }'
  ```

[API usage examples](docs/api_request_examples.md)

[API reference](docs/api.md)

## 📦 Dependencies
- Golang
- OpenAI API (for LLM)
- Ollama LLM (alternative model support)
- Echo framework (for API handling)
- Make (for running build commands)

## 🤝 Contributing
Contributions are welcome! Please follow these steps:
1. Fork the repository
2. Create a new branch (`feature-branch`)
3. Commit changes and push
4. Create a pull request

[Project structure doc](docs/project_structure.md)

## 📜 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Contact
For issues, reach out via GitHub Issues or email: [stefanouorestis@gmail.com](mailto:stefanouorestis@gmail.com)

