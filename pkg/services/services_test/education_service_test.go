package services_test

import (
	"investbot/pkg/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEducationServiceGenerateRagResponse_SuccessWithExistingConversation(t *testing.T) {
	// Arrange
	mockSessionService := new(MockSessionService)
	mockLlm := new(MockLlm)
	service := services.EducationServiceRag{
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
