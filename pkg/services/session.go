package services

type SessionService interface {
	GetConversationBySessionId(sessionId string) ([]map[string]string, error)
}

type MockSessionService struct{}

func (sv MockSessionService) GetConversationBySessionId(sessionId string) ([]map[string]string, error) {
	return []map[string]string{}, nil
}
