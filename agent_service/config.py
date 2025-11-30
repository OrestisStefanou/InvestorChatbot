from enum import Enum

from pydantic_settings import BaseSettings, SettingsConfigDict

class LLMProvider(str, Enum):
    OPENAI = "openai"
    GOOGLE = "google"

class Settings(BaseSettings):
    MONGO_URI: str
    MONGO_DB_NAME: str
    USER_CONTEXT_COLLECTION_NAME: str
    SESSION_COLLECTION_NAME: str
    LLM_PROVIDER: LLMProvider
    LLM_MODEL: str
    OPENAI_API_KEY: str | None = None
    GOOGLE_API_KEY: str | None = None
    MCP_SERVER_URL: str

    model_config = SettingsConfigDict(env_file=".env")

settings = Settings()
