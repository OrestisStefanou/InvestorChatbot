package openAI

import (
	"errors"
	"investbot/pkg/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOpenAiClient is a mock implementation of the OpenAiClient interface
type MockOpenAiClient struct {
	mock.Mock
}

// Chat is a mock method for the Chat function of the OpenAiClient interface
func (m *MockOpenAiClient) Chat(parameters ChatParameters, responseChannel chan<- string) error {
	args := m.Called(parameters, responseChannel)
	return args.Error(0)
}

func TestGenerateResponse(t *testing.T) {
	mockClient := new(MockOpenAiClient)
	llm := OpenAiLLM{
		modelName:   "test-model",
		client:      mockClient,
		temperature: 0.7,
	}

	conversation := []map[string]string{
		{"role": "user", "content": "Hello"},
		{"role": "assistant", "content": "Hi there!"},
	}

	messages := []services.Message{
		{Role: services.User, Content: "Hello"},
		{Role: services.Assistant, Content: "Hi there!"},
	}

	responseChannel := make(chan<- string, 10)
	defer close(responseChannel)

	mockClient.On("Chat", ChatParameters{
		ModelName:   "test-model",
		Temperature: 0.7,
		Messages:    conversation,
	}, responseChannel).Return(nil)

	err := llm.GenerateResponse(messages, responseChannel)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestGenerateResponse_Error(t *testing.T) {
	mockClient := new(MockOpenAiClient)
	llm := OpenAiLLM{
		modelName:   "test-model",
		client:      mockClient,
		temperature: 0.7,
	}

	conversation := []map[string]string{
		{"role": "user", "content": "Hello"},
		{"role": "assistant", "content": "Hi there!"},
	}

	messages := []services.Message{
		{Role: services.User, Content: "Hello"},
		{Role: services.Assistant, Content: "Hi there!"},
	}

	responseChannel := make(chan<- string, 10)
	defer close(responseChannel)

	mockClient.On("Chat", ChatParameters{
		ModelName:   "test-model",
		Temperature: 0.7,
		Messages:    conversation,
	}, responseChannel).Return(errors.New("API error"))

	err := llm.GenerateResponse(messages, responseChannel)
	assert.Error(t, err)
	assert.Equal(t, "API error", err.Error())
	mockClient.AssertExpectations(t)
}
