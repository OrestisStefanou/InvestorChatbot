package services

import (
	"fmt"
	"testing"

	"investbot/pkg/domain"
	"investbot/pkg/services"
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

type MockSessionService struct {
	mock.Mock
}

func (m *MockSessionService) GetConversationBySessionId(sessionId string) ([]map[string]string, error) {
	args := m.Called(sessionId)
	return args.Get(0).([]map[string]string), args.Error(1)
}

type MockLlm struct {
	mock.Mock
}

func (m *MockLlm) GenerateResponse(conversation []map[string]string, responseChannel chan<- string) error {
	args := m.Called(conversation, responseChannel)
	return args.Error(0)
}

func TestGenerateRagResponse_SuccessWithExistingConversation(t *testing.T) {
	// Arrange
	mockDataService := new(MockSectorDataService)
	mockSessionService := new(MockSessionService)
	mockLlm := new(MockLlm)
	service := services.SectorServiceRag{
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

type sectorContext struct {
	sector       domain.Sector
	sectorStocks []domain.SectorStock
}

func TestGenerateRagResponse_SuccessWithNewConversation(t *testing.T) {
	// Arrange
	mockDataService := new(MockSectorDataService)
	mockSessionService := new(MockSessionService)
	mockLlm := new(MockLlm)
	service := services.SectorServiceRag{
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

func TestGenerateRagResponse_ErrorFetchingConversation(t *testing.T) {
	// Arrange
	mockDataService := new(MockSectorDataService)
	mockSessionService := new(MockSessionService)
	mockLlm := new(MockLlm)
	service := services.SectorServiceRag{
		DataService:    mockDataService,
		Llm:            mockLlm,
		SessionService: mockSessionService,
	}

	sessionId := "test-session"
	question := "What are the top stocks in the technology sector?"
	expectedError := "failed to fetch conversation"
	responseChannel := make(chan string, 1)

	mockSessionService.On("GetConversationBySessionId", sessionId).Return([]map[string]string(nil), fmt.Errorf(expectedError))
	// Act
	err := service.GenerateRagResponse(sessionId, question, responseChannel)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockSessionService.AssertExpectations(t)
}

func TestGenerateRagResponse_ErrorFetchingSectors(t *testing.T) {
	// Arrange
	mockDataService := new(MockSectorDataService)
	mockSessionService := new(MockSessionService)
	mockLlm := new(MockLlm)
	service := services.SectorServiceRag{
		DataService:    mockDataService,
		Llm:            mockLlm,
		SessionService: mockSessionService,
	}

	sessionId := "test-session"
	question := "What are the top stocks in the technology sector?"
	newConversation := []map[string]string{}
	expectedError := "failed to fetch sectors"
	responseChannel := make(chan string, 1)

	mockSessionService.On("GetConversationBySessionId", sessionId).Return(newConversation, nil)
	mockDataService.On("GetSectors").Return([]domain.Sector(nil), fmt.Errorf(expectedError))

	// Act
	err := service.GenerateRagResponse(sessionId, question, responseChannel)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockSessionService.AssertExpectations(t)
	mockDataService.AssertExpectations(t)
}

func TestGenerateRagResponse_ErrorFetchingSectorStocks(t *testing.T) {
	// Arrange
	mockDataService := new(MockSectorDataService)
	mockSessionService := new(MockSessionService)
	mockLlm := new(MockLlm)
	service := services.SectorServiceRag{
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
	expectedError := "failed to fetch sector stocks"
	responseChannel := make(chan string, 1)

	mockSessionService.On("GetConversationBySessionId", sessionId).Return(newConversation, nil)
	mockDataService.On("GetSectors").Return(sectors, nil)
	mockDataService.On("GetSectorStocks", "technology").Return([]domain.SectorStock(nil), fmt.Errorf(expectedError))

	// Act
	err := service.GenerateRagResponse(sessionId, question, responseChannel)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockSessionService.AssertExpectations(t)
	mockDataService.AssertExpectations(t)
}

func TestGenerateRagResponse_ErrorGeneratingResponse(t *testing.T) {
	// Arrange
	mockDataService := new(MockSectorDataService)
	mockSessionService := new(MockSessionService)
	mockLlm := new(MockLlm)
	service := services.SectorServiceRag{
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
	expectedError := "LLM response generation failed"
	responseChannel := make(chan<- string, 1)

	mockSessionService.On("GetConversationBySessionId", sessionId).Return(existingConversation, nil)
	mockDataService.On("GetSectors").Return(sectors, nil)
	mockDataService.On("GetSectorStocks", "technology").Return(sectorStocks, nil)
	mockLlm.On(
		"GenerateResponse",
		append(existingConversation, map[string]string{"role": "user", "content": question}),
		responseChannel,
	).Return(fmt.Errorf(expectedError))

	// Act
	err := service.GenerateRagResponse(sessionId, question, responseChannel)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError)
	mockSessionService.AssertExpectations(t)
	mockLlm.AssertExpectations(t)
}
