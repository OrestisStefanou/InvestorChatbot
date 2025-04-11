# Investor Chatbot

![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)
![GitHub issues](https://img.shields.io/github/issues/yourusername/repo-name)
![GitHub stars](https://img.shields.io/github/stars/yourusername/repo-name?style=social)

## 📌 Overview
Investor Chatbot is a conversational AI designed to help beginner investors make informed investment decisions based on financial data, industry trends, and market news. It can answer questions and provide insights into:

- General investing principles
- Stock sectors and industries
- Stock financials (income statement, cash flows, balance sheets)
- Exchange-Traded Funds (ETFs)
- Market and stock news

## 🚀 Features
- AI-powered chatbot with contextual understanding of investment topics
- Real-time stock and market data analysis
- Financial statement insights (Income Statement, Balance Sheet, Cash Flow)
- Sector and industry trends
- ETF recommendations and analysis
- Market news aggregation

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
   OPEN_API_KEY=your_api_key_here
   ```
5. Run the chatbot:
   ```bash
   make run_investbot
   ```

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

## 📜 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Contact
For issues, reach out via GitHub Issues or email: [stefanouorestis@gmail.com](mailto:stefanouorestis@gmail.com)

