package config

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
	return Config{
		OpenAiKey:     "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",
		OpenAiBaseUrl: "https://api.openai.com/v1",

		OllamaBaseUrl: "http://localhost:11434",

		FaqLimit:     5,
		ConvMsgLimit: 10,
	}, nil
}
