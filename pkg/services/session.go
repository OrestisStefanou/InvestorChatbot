package services

import (
	"fmt"
	"investbot/pkg/errors"
	"sync"

	"github.com/google/uuid"
)

type SessionService interface {
	GetConversationBySessionId(sessionId string) ([]Message, error)
	SetConversationForSessionId(conversation []Message, sessionId string) error
	CreateNewSession() (sessionId string, err error)
	AddMessage(sessionId string, msg Message) error
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
	rwMutex      sync.RWMutex
	sessions     map[string][]Message
	convMsgLimit int
}

func NewInMemorySession(convMsgLimit int) (*InMemorySession, error) {
	sessions := make(map[string][]Message)
	return &InMemorySession{sessions: sessions, convMsgLimit: convMsgLimit}, nil
}

func (s *InMemorySession) GetConversationBySessionId(sessionId string) ([]Message, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	conversation, ok := s.sessions[sessionId]
	if !ok {
		return nil, errors.SessionNotFoundError{Message: fmt.Sprintf("sessionID: %s not found", sessionId)}
	}

	limit := s.convMsgLimit
	if limit > len(conversation) {
		limit = len(conversation) // Avoid out-of-bounds error
	}
	MostRecentConv := conversation[len(conversation)-limit:]

	copySlice := make([]Message, len(MostRecentConv))

	// Copy elements from original to new slice
	copy(copySlice, MostRecentConv)

	return copySlice, nil
}

func (s *InMemorySession) AddMessage(sessionId string, msg Message) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	conversation, ok := s.sessions[sessionId]
	if !ok {
		return errors.SessionNotFoundError{Message: fmt.Sprintf("sessionID: %s not found", sessionId)}
	}

	conversation = append(conversation, msg)
	s.sessions[sessionId] = conversation

	return nil
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
