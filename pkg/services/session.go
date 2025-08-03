package services

import (
	"fmt"
	"investbot/pkg/errors"
	"sync"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SessionService interface {
	GetConversationBySessionId(sessionId string) ([]Message, error)
	CreateNewSession() (sessionId string, err error)
	AddMessage(sessionId string, msg Message) error
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

func (s *InMemorySession) CreateNewSession() (sessionId string, err error) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	sessionId = uuid.NewString()
	s.sessions[sessionId] = []Message{}
	return
}

type MongoDBSession struct {
	client         *mongo.Client
	uri            string
	dbName         string
	collectionName string
}

func NewMongoDBSession(client *mongo.Client, uri, dbName, collectionName string) (*MongoDBSession, error) {
	return &MongoDBSession{
		client:         client,
		uri:            uri,
		dbName:         dbName,
		collectionName: collectionName,
	}, nil
}

func (s *MongoDBSession) GetConversationBySessionId(sessionId string) ([]Message, error) {
	// todo
	return nil, nil
}

func (s *MongoDBSession) CreateNewSession() (sessionId string, err error) {
	// todo
	return "", nil
}

func (s *MongoDBSession) AddMessage(sessionId string, msg Message) error {
	// todo
	return nil
}
