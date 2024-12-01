package services

import (
	"fmt"
	"testing"

	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSectorDataService struct {
	mock.Mock
}

func (m *MockSectorDataService) GetSectorStocks(sector string) ([]domain.SectorStock, error) {
	args := m.Called(sector)
	return args.Get(0).([]domain.SectorStock), args.Error(1)
}

func (m *MockSectorDataService) GetSectors() ([]domain.Sector, error) {
	args := m.Called()
	return args.Get(0).([]domain.Sector), args.Error(1)
}

func TestSectorServiceGenerateRagResponse_SuccessWithExistingConversation(t *testing.T) {
	// Arrange
	mockDataService := new(MockSectorDataService)
	mockSessionService := new(SessionServiceMock)
	mockLlm := new(MockLlm)
	service := SectorServiceRag{
		DataService:    mockDataService,
		Llm:            mockLlm,
		SessionService: mockSessionService,
	}

	sessionId := "test-session"
	question := "What are the top stocks in the technology sector?"
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

func TestSectorServiceGenerateRagResponse_SuccessWithNewConversation(t *testing.T) {
	// Arrange
	mockDataService := new(MockSectorDataService)
	mockSessionService := new(SessionServiceMock)
	mockLlm := new(MockLlm)
	service := SectorServiceRag{
		DataService:    mockDataService,
		Llm:            mockLlm,
		SessionService: mockSessionService,
	}

	sessionId := "test-session"
	question := "What are the top stocks in the technology sector?"
	newConversation := []map[string]string{}
	sectors := []domain.Sector{
		{UrlName: "technology", Name: "Technology"},
		{UrlName: "finance", Name: "Finance"},
	}
	sectorStocks := []domain.SectorStock{
		{CompanyName: "TechCorp"},
		{CompanyName: "Innovatech"},
		{CompanyName: "CyberDyn"},
		{CompanyName: "NextGen"},
		{CompanyName: "QuantumAI"},
		{CompanyName: "ExtraStock"},
	}
	var expectedContext string
	for i := 0; i < len(sectors); i++ {
		// Keep only the top 5 stocks for each sector
		context := sectorContext{
			sector:       sectors[i],
			sectorStocks: sectorStocks[:5],
		}
		expectedContext += fmt.Sprintf("%+v\n", context)
	}

	expectedPrompt := fmt.Sprintf(prompts.SectorsPrompt, expectedContext)
	responseChannel := make(chan<- string, 1)

	mockSessionService.On("GetConversationBySessionId", sessionId).Return(newConversation, nil)
	mockDataService.On("GetSectors").Return(sectors, nil)
	mockDataService.On("GetSectorStocks", "technology").Return(sectorStocks, nil)
	mockDataService.On("GetSectorStocks", "finance").Return(sectorStocks, nil)
	mockLlm.On("GenerateResponse", []map[string]string{
		{"role": "system", "content": expectedPrompt},
		{"role": "user", "content": question},
	}, responseChannel).Return(nil)

	// Act
	err := service.GenerateRagResponse(sessionId, question, responseChannel)

	// Assert
	assert.NoError(t, err)
	mockSessionService.AssertExpectations(t)
	mockDataService.AssertExpectations(t)
	mockLlm.AssertExpectations(t)
}

func TestSectorServiceGenerateRagResponse_ErrorFetchingConversation(t *testing.T) {
	// Arrange
	mockDataService := new(MockSectorDataService)
	mockSessionService := new(SessionServiceMock)
	mockLlm := new(MockLlm)
	service := SectorServiceRag{
		DataService:    mockDataService,
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

func TestSectorServiceGenerateRagResponse_ErrorFetchingSectors(t *testing.T) {
	// Arrange
	mockDataService := new(MockSectorDataService)
	mockSessionService := new(SessionServiceMock)
	mockLlm := new(MockLlm)
	service := SectorServiceRag{
		DataService:    mockDataService,
		Llm:            mockLlm,
		SessionService: mockSessionService,
	}

	sessionId := "test-session"
	question := "What are the top stocks in the technology sector?"
	newConversation := []map[string]string{}
	expectedErrorMessage := "failed to fetch sectors"
	expectedError := DataServiceError{
		Message: fmt.Sprintf("GetSectors failed: %s", expectedErrorMessage),
	}
	responseChannel := make(chan string, 1)

	mockSessionService.On("GetConversationBySessionId", sessionId).Return(newConversation, nil)
	mockDataService.On("GetSectors").Return([]domain.Sector(nil), fmt.Errorf(expectedErrorMessage))

	// Act
	err := service.GenerateRagResponse(sessionId, question, responseChannel)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	mockSessionService.AssertExpectations(t)
	mockDataService.AssertExpectations(t)
}

func TestSectorServiceGenerateRagResponse_ErrorFetchingSectorStocks(t *testing.T) {
	// Arrange
	mockDataService := new(MockSectorDataService)
	mockSessionService := new(SessionServiceMock)
	mockLlm := new(MockLlm)
	service := SectorServiceRag{
		DataService:    mockDataService,
		Llm:            mockLlm,
		SessionService: mockSessionService,
	}

	sessionId := "test-session"
	question := "What are the top stocks in the technology sector?"
	newConversation := []map[string]string{}
	sectors := []domain.Sector{
		{UrlName: "technology", Name: "Technology"},
	}
	expectedErrorMessage := "failed to fetch sector stocks"
	expectedError := DataServiceError{
		Message: fmt.Sprintf("GetSectorStocks failed: %s", expectedErrorMessage),
	}
	responseChannel := make(chan string, 1)

	mockSessionService.On("GetConversationBySessionId", sessionId).Return(newConversation, nil)
	mockDataService.On("GetSectors").Return(sectors, nil)
	mockDataService.On("GetSectorStocks", "technology").Return([]domain.SectorStock(nil), fmt.Errorf(expectedErrorMessage))

	// Act
	err := service.GenerateRagResponse(sessionId, question, responseChannel)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	mockSessionService.AssertExpectations(t)
	mockDataService.AssertExpectations(t)
}

func TestSectorServiceGenerateRagResponse_ErrorGeneratingResponse(t *testing.T) {
	// Arrange
	mockDataService := new(MockSectorDataService)
	mockSessionService := new(SessionServiceMock)
	mockLlm := new(MockLlm)
	service := SectorServiceRag{
		DataService:    mockDataService,
		Llm:            mockLlm,
		SessionService: mockSessionService,
	}

	sessionId := "test-session"
	question := "What are the top stocks in the technology sector?"
	existingConversation := []map[string]string{
		{"role": "system", "content": "Welcome to InvestBot."},
	}
	sectors := []domain.Sector{
		{UrlName: "technology", Name: "Technology"},
	}
	sectorStocks := []domain.SectorStock{
		{CompanyName: "TechCorp"},
		{CompanyName: "Innovatech"},
		{CompanyName: "CyberDyn"},
		{CompanyName: "NextGen"},
		{CompanyName: "QuantumAI"},
	}
	expectedErrorMessage := "LLM response generation failed"
	expectedError := RagError{
		Message: fmt.Sprintf("GenerateResponse failed: %s", expectedErrorMessage),
	}
	responseChannel := make(chan<- string, 1)

	mockSessionService.On("GetConversationBySessionId", sessionId).Return(existingConversation, nil)
	mockDataService.On("GetSectors").Return(sectors, nil)
	mockDataService.On("GetSectorStocks", "technology").Return(sectorStocks, nil)
	mockLlm.On(
		"GenerateResponse",
		append(existingConversation, map[string]string{"role": "user", "content": question}),
		responseChannel,
	).Return(fmt.Errorf(expectedErrorMessage))

	// Act
	err := service.GenerateRagResponse(sessionId, question, responseChannel)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	mockSessionService.AssertExpectations(t)
	mockLlm.AssertExpectations(t)
}
