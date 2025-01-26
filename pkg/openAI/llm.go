package openAI

import (
	"fmt"
	"investbot/pkg/services"
)

type OpenAiClientInterface interface {
	Chat(parameters ChatParameters, responseChannel chan<- string) error
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
// - conversation: A slice of Message
func (llm OpenAiLLM) GenerateResponse(conversation []services.Message, responseChannel chan<- string) error {
	// Send the messages to the OpenAI API
	messages := make([]map[string]string, 0, len(conversation))
	for _, m := range conversation {
		msg := make(map[string]string)
		msg["role"] = string(m.Role)
		msg["content"] = m.Content
		messages = append(messages, msg)
	}
	parameters := ChatParameters{
		ModelName:   string(llm.modelName),
		Temperature: llm.temperature,
		Messages:    messages,
	}
	if err := llm.client.Chat(parameters, responseChannel); err != nil {
		return err
	}
	return nil
}
