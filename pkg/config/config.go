package config

import (
	"investbot/pkg/openAI"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type LlmProvider string

const (
	OPEN_AI LlmProvider = "OPEN_AI"
	OLLAMA              = "OLLAMA"
)

type Config struct {
	// OpenAI configs
	OpenAiKey       string
	OpenAiBaseUrl   string
	OpenAiModelName openAI.ModelName

	// Ollama configs
	OllamaBaseUrl   string
	OllamaModelName string

	// App configs
	LlmProvider          LlmProvider // Valid values are: "OPEN_AI", "OLLAMA"
	FaqLimit             int         // Number of faq to return in through the endpoint
	ConvMsgLimit         int         // The number of most recent messages to get from a session
	BaseLlmTemperature   float32     // The temperature to use for the base llm(currently there is only one llm that is used in all the rags)
	FollowUpQuestionsNum int         // The number of follow-up questions that the GET /follow_up_questions will return
	CacheTtl             int         // The ttl for the cache in seconds

	// Badger configs
	BadgerDbPath string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load() // Load .env file, ignore errors if not found
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

	llmProvider := getEnv("LLM_PROVIDER", "OPEN_AI")

	openAiModelName := getEnv("OPEN_AI_MODEL_NAME", "gpt-4o-mini")

	ollamaModelName := getEnv("OLLAMA_MODEL_NAME", "llama3.2")

	return Config{
		OpenAiKey:            getEnv("OPEN_AI_API_KEY", ""),
		OpenAiBaseUrl:        getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		OllamaBaseUrl:        getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
		FaqLimit:             faqLimit,
		ConvMsgLimit:         convMsgLimit,
		LlmProvider:          LlmProvider(llmProvider),
		OpenAiModelName:      openAI.ModelName(openAiModelName),
		OllamaModelName:      ollamaModelName,
		BaseLlmTemperature:   getEnvFloat32("BASE_LLM_TEMPERATURE", 0.2),
		FollowUpQuestionsNum: followUpQuestionsNum,
		CacheTtl:             cacheTtl,
		BadgerDbPath:         getEnv("BADGER_DB_PATH", "badger.db"),
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
