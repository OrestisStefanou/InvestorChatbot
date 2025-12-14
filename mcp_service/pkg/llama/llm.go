package llama

import "investbot/pkg/services"

type LlamaClientInterface interface {
	Chat(parameters ChatParameters, responseChannel chan<- string) error
}

type ModelName string

const (
	LLAMA_3_2 ModelName = "llama3.2"
)

type LllamaLLM struct {
	modelName   ModelName
	client      LlamaClientInterface
	temperature float32
}

func NewLlamaLLM(modelName ModelName, client LlamaClientInterface, temperature float32) (*LllamaLLM, error) {
	return &LllamaLLM{
		modelName:   modelName,
		client:      client,
		temperature: temperature,
	}, nil
}

func (llm LllamaLLM) GenerateResponse(conversation []services.Message, responseChannel chan<- string) error {
	messages := make([]Message, 0, len(conversation))
	for _, m := range conversation {
		msg := Message{
			Role:    string(m.Role),
			Content: m.Content,
		}
		messages = append(messages, msg)
	}

	parameters := ChatParameters{
		ModelName: string(llm.modelName),
		Messages:  messages,
		Options: Options{
			Temperature: llm.temperature,
		},
		Stream: true,
	}
	if err := llm.client.Chat(parameters, responseChannel); err != nil {
		return err
	}
	return nil
}

func (llm LllamaLLM) GetLlmName() string {
	return string(llm.modelName)
}
