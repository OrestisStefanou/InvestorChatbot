package services

type SessionService interface {
	GetConversationBySessionId(sessionId string) ([]map[string]string, error)
}
