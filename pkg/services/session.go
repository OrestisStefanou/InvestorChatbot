package services

type SessionService interface {
	GetConversationBySessionId(sessionId string) ([]Message, error)
}

type MockSessionService struct{}

func (sv MockSessionService) GetConversationBySessionId(sessionId string) ([]Message, error) {
	return []Message{}, nil
}
