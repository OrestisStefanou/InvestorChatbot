package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// OpenAI configs
	OpenAiKey     string
	OpenAiBaseUrl string

	// Ollama configs
	OllamaBaseUrl string

	// App configs
	FaqLimit     int // Number of faq to return in through the endpoint
	ConvMsgLimit int // The number of most recent messages to get from a session
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load() // Load .env file, ignore errors if not found

	faqLimit, err := strconv.Atoi(getEnv("FAQ_LIMIT", "5"))
	if err != nil {
		faqLimit = 5
	}

	convMsgLimit, err := strconv.Atoi(getEnv("CONV_MSG_LIMIT", "10"))
	if err != nil {
		convMsgLimit = 10
	}

	return Config{
		OpenAiKey:     getEnv("OPEN_AI_API_KEY", ""),
		OpenAiBaseUrl: getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		OllamaBaseUrl: getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
		FaqLimit:      faqLimit,
		ConvMsgLimit:  convMsgLimit,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
