package services

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type SessionService interface {
	GetConversationBySessionId(sessionId string) ([]Message, error)
	SetConversationForSessionId(conversation []Message, sessionId string) error
	CreateNewSession() (sessionId string, err error)
}

type MockSessionService struct{}

func (sv MockSessionService) GetConversationBySessionId(sessionId string) ([]Message, error) {
	return []Message{}, nil
}

func (sv MockSessionService) SetConversationForSessionId(conversation []Message, sessionId string) error {
	return nil
}

func (sv MockSessionService) CreateNewSession() (sessionId string, err error) {
	return "mock_session_id", nil
}

type InMemorySession struct {
	rwMutex  sync.RWMutex
	sessions map[string][]Message
}

func NewInMemorySession() (*InMemorySession, error) {
	sessions := make(map[string][]Message)
	return &InMemorySession{sessions: sessions}, nil
}

func (s *InMemorySession) GetConversationBySessionId(sessionId string) ([]Message, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	conversation, ok := s.sessions[sessionId]
	if !ok {
		return nil, fmt.Errorf("Conversation with sessionID: %s not found", sessionId)
	}

	return conversation, nil
}

func (s *InMemorySession) SetConversationForSessionId(conversation []Message, sessionId string) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	s.sessions[sessionId] = conversation
	return nil
}

func (s *InMemorySession) CreateNewSession() (sessionId string, err error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	sessionId = uuid.NewString()
	s.sessions[sessionId] = []Message{}
	return
}
