# 📁 Project Structure

This project is organized using a modular structure to improve readability, maintainability, and scalability. Below is an overview of each major directory and its purpose.

---

## 🗂 Root Files

- **Makefile** – Contains useful commands for building, running, or testing the application.
- **go.mod / go.sum** – Go module dependencies.
- **README.md** – Main documentation file.
- **docs/** – Additional documentation (e.g., API and configuration).

---

## 🧠 `pkg/` – Core Application Logic

This directory contains all the core packages of the app.

### 🔹 `pkg/domain/`
Defines the **domain models** used throughout the application, including:

- Stocks
- ETFs
- Sectors
- Industries
- News
- Financials
- Super investors

### 🔹 `pkg/errors/`
Houses all **custom error definitions** used across the app to provide descriptive and structured error responses.

### 🔹 `pkg/handlers/`
Contains the **HTTP endpoint handlers** that process incoming API requests and return responses. These handlers map directly to your route definitions.

### 🔹 `pkg/llama/`
Includes the integration logic for interacting with **Ollama LLMs**.

### 🔹 `pkg/marketDataScraper/`
Implements logic to **scrape and fetch market data** such as:

- Financial statements
- Forecasts
- ETFs and sectors
- Industry stocks
- News
- Super investor portfolios

Also includes sample responses under `example_responses/` for development and testing.

### 🔹 `pkg/openAI/`
Handles interaction with **OpenAI’s API**, including client setup and model usage.

### 🔹 `pkg/repositories/`
Reserved for any **database interaction or persistence layer** code.

### 🔹 `pkg/services/`
Encapsulates the **business logic** of the application.

Key subdirectories and files:

- `faq/`: Logic for handling FAQs (balance sheets, income statements, etc.)
- `prompts/`: Prompt templates used in LLM-based features
- `session.go`, `chat_service.go`, etc.: Higher-level logic driving feature behavior

### 🔹 `pkg/config/`
The **configuration loader** for environment variables and `.env` file settings.

### 🔹 `pkg/utils/`
Houses general-purpose **utility functions**.

---

## 🚀 `cmd/` – Application Entry Points

Each folder under `cmd/` corresponds to an executable application or service.

- **investbot/**: Main entry point for the core InvestBot application
- **temp/**: Temporary or experimental logic

Each contains a `main.go` file as the program entry point.

---

## 🗃 `docs/` – Documentation

- **api.md**: API endpoint documentation
- **config.md**: Configuration/environment variable reference

---

## ✅ Summary

| Folder                | Purpose                                        |
|-----------------------|------------------------------------------------|
| `pkg/domain`          | Domain models                                  |
| `pkg/errors`          | Custom error types                             |
| `pkg/handlers`        | API endpoint handlers                          |
| `pkg/llama`           | Ollama integration                             |
| `pkg/marketDataScraper`| Market data scraping logic                    |
| `pkg/repositories`    | Database layer                                 |
| `pkg/services`        | Business logic                                 |
| `pkg/config`          | Configuration loading                          |
| `pkg/utils`           | Shared utility functions                       |
| `cmd/`                | Application entry points                       |
| `docs/`               | Documentation files                            |

---
