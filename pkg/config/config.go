package config

import (
	"investbot/pkg/gemini"
	"investbot/pkg/openAI"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type LlmProvider string

const (
	OPEN_AI LlmProvider = "OPEN_AI"
	OLLAMA  LlmProvider = "OLLAMA"
	GEMINI  LlmProvider = "GEMINI"
)

type DatabaseProvider string

const (
	MONGO_DB  DatabaseProvider = "MONGO_DB"
	BADGER_DB DatabaseProvider = "BADGER"
)

type SessionStorageProvider string

const (
	MONGO_DB_STORAGE  SessionStorageProvider = "MONGO_DB"
	IN_MEMORY_STORAGE SessionStorageProvider = "MEMORY"
)

type MongoDBConfig struct {
	Uri                        string
	DBName                     string
	SessionCollectionName      string
	UserContextColletionName   string
	TopicAndTagsCollectionName string
	RagResponsesCollectionName string
}

type Config struct {
	// OpenAI configs
	OpenAiKey       string
	OpenAiBaseUrl   string
	OpenAiModelName openAI.ModelName

	// Ollama configs
	OllamaBaseUrl   string
	OllamaModelName string

	// Gemini configs
	GeminiKey       string
	GeminiModelName gemini.ModelName

	// App configs
	LlmProvider            LlmProvider // Valid values are: "OPEN_AI", "OLLAMA"
	FaqLimit               int         // Number of faq to return in through the endpoint
	ConvMsgLimit           int         // The number of most recent messages to get from a session
	BaseLlmTemperature     float32     // The temperature to use for the base llm(currently there is only one llm that is used in all the rags)
	FollowUpQuestionsNum   int         // The number of follow-up questions that the GET /follow_up_questions will return
	CacheTtl               int         // The ttl for the cache in seconds
	DatabaseProvider       DatabaseProvider
	SessionStorageProvider SessionStorageProvider

	// Badger configs
	BadgerDbPath string

	// MongoDB configs
	MongoDBConf MongoDBConfig
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	faqLimit, err := strconv.Atoi(getEnv("FAQ_LIMIT", "5"))
	if err != nil {
		faqLimit = 5
	}

	convMsgLimit, err := strconv.Atoi(getEnv("CONV_MSG_LIMIT", "10"))
	if err != nil {
		convMsgLimit = 10
	}

	followUpQuestionsNum, err := strconv.Atoi(getEnv("FOLLOW_UP_QUESTIONS_NUM", "5"))
	if err != nil {
		followUpQuestionsNum = 5
	}

	cacheTtl, err := strconv.Atoi(getEnv("CACHE_TTL", "3600"))
	if err != nil {
		cacheTtl = 3600
	}

	dbProvider := getEnv("DATABASE_PROVIDER", "BADGER")

	sessionStorage := getEnv("SESSION_STORAGE_PROVIDER", "MEMORY")

	llmProvider := getEnv("LLM_PROVIDER", "OPEN_AI")

	openAiModelName := getEnv("OPEN_AI_MODEL_NAME", "gpt-4o-mini")

	ollamaModelName := getEnv("OLLAMA_MODEL_NAME", "llama3.2")

	geminiModelName := getEnv("GEMINI_MODEL_NAME", "gemini-2.0-flash")

	return Config{
		OpenAiKey:            getEnv("OPEN_AI_API_KEY", ""),
		OpenAiBaseUrl:        getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		OllamaBaseUrl:        getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
		GeminiKey:            getEnv("GEMINI_API_KEY", ""),
		FaqLimit:             faqLimit,
		ConvMsgLimit:         convMsgLimit,
		LlmProvider:          LlmProvider(llmProvider),
		OpenAiModelName:      openAI.ModelName(openAiModelName),
		GeminiModelName:      gemini.ModelName(geminiModelName),
		OllamaModelName:      ollamaModelName,
		BaseLlmTemperature:   getEnvFloat32("BASE_LLM_TEMPERATURE", 0.2),
		FollowUpQuestionsNum: followUpQuestionsNum,
		CacheTtl:             cacheTtl,
		BadgerDbPath:         getEnv("BADGER_DB_PATH", "badger.db"),
		MongoDBConf: MongoDBConfig{
			Uri:                        getEnv("MONGO_DB_URI", ""),
			DBName:                     getEnv("MONGO_DB_NAME", ""),
			SessionCollectionName:      getEnv("MONGO_DB_SESSION_COLLECTION_NAME", "session"),
			UserContextColletionName:   getEnv("MONGO_DB_USER_CONTEXT_COLLECTION_NAME", "user_context"),
			TopicAndTagsCollectionName: getEnv("MONGO_DB_TOPIC_AND_TAGS_COLLECTION_NAME", "topic_and_tags"),
			RagResponsesCollectionName: getEnv("MONGO_DB_RAG_RESPONSES_COLLECTION_NAME", "rag_responses"),
		},
		DatabaseProvider:       DatabaseProvider(dbProvider),
		SessionStorageProvider: SessionStorageProvider(sessionStorage),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvFloat32(key string, fallback float32) float32 {
	if value, exists := os.LookupEnv(key); exists {
		if floatValue, err := strconv.ParseFloat(value, 32); err == nil {
			return float32(floatValue)
		}
	}
	return fallback
}
