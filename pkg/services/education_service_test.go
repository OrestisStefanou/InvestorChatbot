package services

import (
	"fmt"
	"investbot/pkg/services/prompts"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEducationServiceGenerateRagResponse_SuccessWithExistingConversation(t *testing.T) {
	// Arrange
	mockSessionService := new(SessionServiceMock)
	mockLlm := new(MockLlm)
	service := EducationServiceRag{
		Llm:            mockLlm,
		SessionService: mockSessionService,
	}

	sessionId := "test-session"
	question := "What are stocks?"
	existingConversation := []map[string]string{
		{"role": "system", "content": "Welcome to InvestBot."},
	}
	responseChannel := make(chan<- string, 1)

	mockSessionService.On("GetConversationBySessionId", sessionId).Return(existingConversation, nil)
	mockLlm.On("GenerateResponse", append(existingConversation, map[string]string{"role": "user", "content": question}), responseChannel).Return(nil)

	// Act
	err := service.GenerateRagResponse(sessionId, question, responseChannel)

	// Assert
	assert.NoError(t, err)
	mockSessionService.AssertExpectations(t)
	mockLlm.AssertExpectations(t)
}

func TestEducationServiceGenerateRagResponse_SuccessWithNewConversation(t *testing.T) {
	// Arrange
	mockSessionService := new(SessionServiceMock)
	mockLlm := new(MockLlm)
	service := EducationServiceRag{
		Llm:            mockLlm,
		SessionService: mockSessionService,
	}

	sessionId := "test-session"
	question := "What are the top stocks in the technology sector?"
	newConversation := []map[string]string{}

	responseChannel := make(chan<- string, 1)

	mockSessionService.On("GetConversationBySessionId", sessionId).Return(newConversation, nil)
	mockLlm.On("GenerateResponse", []map[string]string{
		{"role": "system", "content": prompts.EducationPrompt},
		{"role": "user", "content": question},
	}, responseChannel).Return(nil)

	// Act
	err := service.GenerateRagResponse(sessionId, question, responseChannel)

	// Assert
	assert.NoError(t, err)
	mockSessionService.AssertExpectations(t)
	mockLlm.AssertExpectations(t)
}

func TestEducationServiceGenerateRagResponse_ErrorFetchingConversation(t *testing.T) {
	// Arrange
	mockSessionService := new(SessionServiceMock)
	mockLlm := new(MockLlm)
	service := EducationServiceRag{
		Llm:            mockLlm,
		SessionService: mockSessionService,
	}

	sessionId := "test-session"
	question := "What are the top stocks in the technology sector?"
	expectedErrorMessage := "failed to fetch conversation"
	expectedError := SessionServiceError{
		Message: fmt.Sprintf("GetConversationBySessionId failed: %s", expectedErrorMessage),
	}
	responseChannel := make(chan string, 1)

	mockSessionService.On("GetConversationBySessionId", sessionId).Return([]map[string]string(nil), fmt.Errorf(expectedErrorMessage))
	// Act
	err := service.GenerateRagResponse(sessionId, question, responseChannel)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	mockSessionService.AssertExpectations(t)
}
