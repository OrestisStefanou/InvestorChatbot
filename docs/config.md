# 🔧 App Configuration Guide

This document explains how to configure and use environment-based settings for the application using the `config` package. It supports flexible setup through a `.env` file or environment variables.
---

## Providers

### LlmProvider
Specifies the Large Language Model (LLM) provider.

- **Type:** `string`
- **Possible values:**
  - `OPEN_AI`
  - `OLLAMA`
  - `GEMINI`

### DatabaseProvider
Specifies the database provider.

- **Type:** `string`
- **Possible values:**
  - `MONGO_DB`
  - `BADGER`

### SessionStorageProvider
Specifies how session data is stored.

- **Type:** `string`
- **Possible values:**
  - `MONGO_DB`
  - `MEMORY`

---

## Main Config Structure

### `Config` Fields

#### OpenAI Configuration
- `OpenAiKey` – API key for OpenAI.  
- `OpenAiBaseUrl` – Base URL for OpenAI API. Default: `https://api.openai.com/v1`
- `OpenAiModelName` – Model name (e.g., `gpt-4o-mini`).

#### Ollama Configuration
- `OllamaBaseUrl` – Base URL for Ollama API. Default: `http://localhost:11434`
- `OllamaModelName` – Model name (e.g., `llama3.2`).

#### Gemini Configuration
- `GeminiKey` – API key for Gemini.
- `GeminiModelName` – Model name (e.g., `gemini-2.0-flash`).

---

### Application Configs
- `LlmProvider` – LLM provider to use. Default: `OPEN_AI`
- `FaqLimit` – Number of FAQs returned by endpoints. Default: `5`
- `ConvMsgLimit` – Number of recent session messages to retrieve. Default: `10`
- `BaseLlmTemperature` – Temperature setting for the base LLM. Default: `0.2`
- `FollowUpQuestionsNum` – Number of follow-up questions to return. Default: `5`
- `CacheTtl` – Cache TTL in seconds. Default: `3600`
- `DatabaseProvider` – Database provider (`MONGO_DB` or `BADGER`).
- `SessionStorageProvider` – Session storage provider (`MONGO_DB` or `MEMORY`).

---

## MongoDB Configuration

### `MongoDBConfig`
- `Uri` – MongoDB connection URI.
- `DBName` – Database name.
- `SessionCollectionName` – Collection for session data. Default: `session`
- `UserContextColletionName` – Collection for user context. Default: `user_context`
- `TopicAndTagsCollectionName` – Collection for topics and tags. Default: `topic_and_tags`
- `RagResponsesCollectionName` – Collection for RAG responses. Default: `rag_responses`

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `OPEN_AI_API_KEY` | `""` | OpenAI API key |
| `OPENAI_BASE_URL` | `https://api.openai.com/v1` | OpenAI API base URL |
| `OLLAMA_BASE_URL` | `http://localhost:11434` | Ollama API base URL |
| `GEMINI_API_KEY` | `""` | Gemini API key |
| `FAQ_LIMIT` | `5` | FAQ results limit |
| `CONV_MSG_LIMIT` | `10` | Conversation message limit |
| `FOLLOW_UP_QUESTIONS_NUM` | `5` | Number of follow-up questions |
| `CACHE_TTL` | `3600` | Cache TTL in seconds |
| `BADGER_DB_PATH` | `badger.db` | BadgerDB file path |
| `MONGO_DB_URI` | `""` | MongoDB connection string |
| `MONGO_DB_NAME` | `""` | MongoDB database name |
| `MONGO_DB_SESSION_COLLECTION_NAME` | `session` | Session collection name |
| `MONGO_DB_USER_CONTEXT_COLLECTION_NAME` | `user_context` | User context collection name |
| `MONGO_DB_TOPIC_AND_TAGS_COLLECTION_NAME` | `topic_and_tags` | Topic and tags collection name |
| `MONGO_DB_RAG_RESPONSES_COLLECTION_NAME` | `rag_responses` | RAG responses collection name |

---

## Loading Configuration
The function `LoadConfig()` loads values from `.env` and applies defaults if variables are missing.

