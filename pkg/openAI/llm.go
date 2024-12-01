package openAI

import "fmt"

type OpenAiClientInterface interface {
	Chat(parameters chatParameters, responseChannel chan<- string) error
}

type ModelName string

const (
	GPT3      ModelName = "gpt-3"
	GPT4_MINI ModelName = "gpt-4o-mini"
)

type OpenAiLLM struct {
	modelName   ModelName
	client      OpenAiClientInterface
	temperature float64
}

func NewOpenAiLLM(modelName ModelName, client OpenAiClientInterface, temperature float64) (*OpenAiLLM, error) {
	switch modelName {
	case GPT3, GPT4_MINI:
	default:
		return nil, fmt.Errorf("invalid model name: %s", modelName)
	}

	// Return a new OpenAiLLM instance
	return &OpenAiLLM{
		modelName:   modelName,
		client:      client,
		temperature: temperature,
	}, nil
}

// GenerateResponse generates a response from the OpenAI language model based on the provided conversatoin.
// It sends the conversation messages to the OpenAI API and streams the response in chunks.
// The response chunks are sent over the responseChannel for real-time processing.
// Params:
// - conversation: A slice of maps that must have the following format
//
//	{
//		{"role": "user", "content": "Hey there!"},
//		{"role": "system", "content": "Hello! How can I help you today?"},
//		{"role": "user", "content": "What is a synonym for big?"},
//	}
func (llm OpenAiLLM) GenerateResponse(conversation []map[string]string, responseChannel chan<- string) error {
	// Send the messages to the OpenAI API
	parameters := chatParameters{
		ModelName:   string(llm.modelName),
		Temperature: llm.temperature,
		Messages:    conversation,
	}
	if err := llm.client.Chat(parameters, responseChannel); err != nil {
		return err
	}
	return nil
}
