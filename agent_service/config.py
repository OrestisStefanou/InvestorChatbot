from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    MONGO_URI: str
    MONGO_DB_NAME: str
    USER_CONTEXT_COLLECTION_NAME: str

    model_config = SettingsConfigDict(env_file=".env")

settings = Settings()
