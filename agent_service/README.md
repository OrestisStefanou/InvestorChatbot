# Agent Service

The `agent-service` is the core backend component for the Investor Chatbot. It is a FastAPI application that provides intelligent, personalized investment advice by leveraging Large Language Models (LLMs) and the Model Context Protocol (MCP).

## 🚀 Features

- **Personalized Investment Advisor**: Acts as a professional investment advisor, tailoring advice to the user's specific context.
- **Multi-LLM Support**: Configurable to use OpenAI, Google Gemini, or Anthropic Claude models.
- **Context Awareness**: Persistently stores and retrieves user context and session history using MongoDB to ensure continuity in conversations.
- **MCP Integration**: Uses the Model Context Protocol to seamlessly integrate with external tools and data sources.

## 🛠️ Prerequisites

- **Python**: >= 3.13
- **MongoDB**: A running MongoDB instance (local or cloud).
- **Package Manager**: `uv` (recommended) or `pip`.

## ⚙️ Configuration

1.  **Copy the example environment file:**
    ```bash
    cp .env.example .env
    ```

2.  **Configure `.env`:**
    Update the `.env` file with your specific configuration:
    - `MONGO_URI`: Your MongoDB connection string.
    - `MONGO_DB_NAME`: Database name (default: `investor_chatbot`).
    - `USER_CONTEXT_COLLECTION_NAME`: Collection for user context.
    - **LLM Keys**: Add API keys for your chosen provider (e.g., `OPENAI_API_KEY`, `GOOGLE_API_KEY`, `ANTHROPIC_API_KEY`).
    - **LLM Settings**: Configure `LLM_PROVIDER`, `LLM_MODEL`, and `TEMPERATURE` as needed (check `config.py` for all options).

## 📦 Installation

This project uses `uv` for fast dependency management.

1.  **Install dependencies:**
    ```bash
    uv sync
    ```
    Or with pip:
    ```bash
    pip install .
    ```

## 🏃‍♂️ Running the Service

You can run the service using the provided `Makefile` or directly with `uv`.

**Using Makefile:**
```bash
make run_agent_service
```

**Using uv directly:**
```bash
uv run fastapi dev main.py
```

The service will start on `http://127.0.0.1:8000` by default.

## 📂 Project Structure

- `main.py`: Application entry point and lifespan management.
- `routers/`: API route definitions (chat, session, user_context).
- `services/`: Business logic, including the core `Agent` and `AgentService`.
- `config.py`: Configuration management using Pydantic Settings.
- `dependencies.py`: Dependency injection providers.
- `.env.example`: Template for environment variables.

## 📡 API Endpoints

### Chat

-   **`POST /chat`**: Send a message to the agent and get a response.
    -   **Request Body**:
        ```json
        {
          "session_id": "string",
          "message": "string"
        }
        ```
    -   **Response**: Returns the agent's textual response.

### Session

-   **`POST /session`**: Create a new chat session.
    -   **Request Body**:
        ```json
        {
          "user_id": "string"
        }
        ```
    -   **Response**: Returns the created session details including `session_id`.

-   **`GET /session/{session_id}`**: Retrieve an existing session.
    -   **Response**: Returns the full session history, including all messages exchanged.

### User Context

-   **`POST /user_context`**: Create a new user context profile.
    -   **Request Body**:
        ```json
        {
          "user_id": "string",
          "user_profile": {},
          "user_portfolio": [
            {
              "asset_class": "string",
              "symbol": "string",
              "name": "string",
              "quantity": 0
            }
          ]
        }
        ```

-   **`GET /user_context/{user_id}`**: Retrieve a user's context.

-   **`PUT /user_context`**: Update an existing user context.
    -   **Request Body**: Same as `POST /user_context`.