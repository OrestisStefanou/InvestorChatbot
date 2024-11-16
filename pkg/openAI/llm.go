package openAI

type OpenAiClientInterface interface {
	Chat(parameters ChatParameters, responseChannel chan<- string) error
}

type OpenAiLLM struct {
	ModelName   string
	Client      OpenAiClientInterface
	Temperature float64
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
	parameters := ChatParameters{
		ModelName:   llm.ModelName,
		Temperature: llm.Temperature,
		Messages:    conversation,
	}
	if err := llm.Client.Chat(parameters, responseChannel); err != nil {
		return err
	}
	return nil
}
