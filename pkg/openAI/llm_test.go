package openAI

import (
	"errors"
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
		ModelName:    "test-model",
		Client:       mockClient,
		SystemPrefix: "System message",
		Temperature:  0.7,
	}

	conversation := []map[string]string{
		{"role": "user", "content": "Hello"},
		{"role": "assistant", "content": "Hi there!"},
	}

	responseChannel := make(chan string, 10)
	defer close(responseChannel)

	expectedMessages := []map[string]string{
		{"role": "system", "content": "System message"},
		{"role": "user", "content": "Hello"},
		{"role": "assistant", "content": "Hi there!"},
	}

	mockClient.On("Chat", ChatParameters{
		ModelName:   "test-model",
		Temperature: 0.7,
		Messages:    expectedMessages,
	}, responseChannel).Return(nil)

	err := llm.GenerateResponse(conversation, responseChannel)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestGenerateResponse_Error(t *testing.T) {
	mockClient := new(MockOpenAiClient)
	llm := OpenAiLLM{
		ModelName:    "test-model",
		Client:       mockClient,
		SystemPrefix: "System message",
		Temperature:  0.7,
	}

	conversation := []map[string]string{
		{"role": "user", "content": "Hello"},
		{"role": "assistant", "content": "Hi there!"},
	}

	responseChannel := make(chan string, 10)
	defer close(responseChannel)

	expectedMessages := []map[string]string{
		{"role": "system", "content": "System message"},
		{"role": "user", "content": "Hello"},
		{"role": "assistant", "content": "Hi there!"},
	}

	mockClient.On("Chat", ChatParameters{
		ModelName:   "test-model",
		Temperature: 0.7,
		Messages:    expectedMessages,
	}, responseChannel).Return(errors.New("API error"))

	err := llm.GenerateResponse(conversation, responseChannel)
	assert.Error(t, err)
	assert.Equal(t, "API error", err.Error())
	mockClient.AssertExpectations(t)
}
