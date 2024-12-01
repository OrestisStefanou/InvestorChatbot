package services

import "github.com/stretchr/testify/mock"

type SessionServiceMock struct {
	mock.Mock
}

func (m *SessionServiceMock) GetConversationBySessionId(sessionId string) ([]map[string]string, error) {
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
