package openAI

type OpenAiClientInterface interface {
	Chat(parameters ChatParameters, responseChannel chan string) error
}

type OpenAiLLM struct {
	ModelName    string
	Client       OpenAiClientInterface
	SystemPrefix string
	Temperature  float64
}

// GenerateResponse generates a response from the OpenAI language model based on the provided conversatoin.
// It sends the conversation messages to the OpenAI API and streams the response in chunks.
// The system message is prepended to the conversation messages before sending them to the API.
// The response chunks are sent over the responseChannel for real-time processing.
func (llm OpenAiLLM) GenerateResponse(conversation []map[string]string, responseChannel chan string) error {
	// Preallocate the slice with the required capacity
	messages := make([]map[string]string, 0, len(conversation)+1)
	// Add the system message
	messages = append(messages, map[string]string{"role": "system", "content": llm.SystemPrefix})
	// Add the conversation messages
	messages = append(messages, conversation...)
	// Send the messages to the OpenAI API
	parameters := ChatParameters{
		ModelName:   llm.ModelName,
		Temperature: llm.Temperature,
		Messages:    messages,
	}
	if err := llm.Client.Chat(parameters, responseChannel); err != nil {
		return err
	}
	return nil
}
