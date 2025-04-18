# 🔧 App Configuration Guide

This document explains how to configure and use environment-based settings for the application using the `config` package. It supports flexible setup through a `.env` file or environment variables.

## 🧠 Supported LLM Providers

The app supports two large language model (LLM) providers:

- **OpenAI**
- **Ollama**

The provider is selected via the `LLM_PROVIDER` environment variable.

---

## 📄 Sample `.env` File

```env
# LLM Provider
LLM_PROVIDER=OPEN_AI  # or OLLAMA

# OpenAI Config
OPEN_AI_API_KEY=your-openai-key
OPENAI_BASE_URL=https://api.openai.com/v1
OPEN_AI_MODEL_NAME=gpt-4o-mini

# Ollama Config
OLLAMA_BASE_URL=http://localhost:11434
OLLAMA_MODEL_NAME=llama3.2

# Application Settings
FAQ_LIMIT=5
CONV_MSG_LIMIT=10
BASE_LLM_TEMPERATURE=0.2
FOLLOW_UP_QUESTIONS_NUM=5
```

---

## ⚙️ Environment Variables Reference

| Variable Name               | Description                                                               | Type     | Default                     |
|----------------------------|---------------------------------------------------------------------------|----------|-----------------------------|
| `LLM_PROVIDER`             | Determines which LLM to use: `OPEN_AI` or `OLLAMA`                        | `string` | `OPEN_AI`                   |
| `OPEN_AI_API_KEY`          | API key for OpenAI                                                        | `string` | `""` (empty)                |
| `OPENAI_BASE_URL`          | Base URL for OpenAI API                                                   | `string` | `https://api.openai.com/v1` |
| `OPEN_AI_MODEL_NAME`       | Model name used with OpenAI                                               | `string` | `gpt-4o-mini`               |
| `OLLAMA_BASE_URL`          | Base URL for Ollama server                                                | `string` | `http://localhost:11434`    |
| `OLLAMA_MODEL_NAME`        | Model name used with Ollama                                               | `string` | `llama3.2`                  |
| `FAQ_LIMIT`                | Number of FAQ items returned by the FAQ endpoint                          | `int`    | `5`                         |
| `CONV_MSG_LIMIT`           | Number of most recent conversation messages to retrieve                   | `int`    | `10`                        |
| `BASE_LLM_TEMPERATURE`     | Temperature setting for the LLM, controls randomness                      | `float`  | `0.2`                       |
| `FOLLOW_UP_QUESTIONS_NUM` | Number of follow-up questions to generate                                 | `int`    | `5`                         |

---

## 🛠 Loading Config in Code

The `LoadConfig` function loads and validates all configuration values:

```go
config, err := config.LoadConfig()
if err != nil {
	log.Fatalf("Error loading config: %v", err)
}
```

It automatically falls back to default values if environment variables are missing or invalid.

---

## 📦 Dependencies

- [github.com/joho/godotenv](https://github.com/joho/godotenv) – for loading environment variables from `.env` files.

---
