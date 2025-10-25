package services

import (
	"context"
	"fmt"
	"investbot/pkg/errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

type MongoDBSessionServiceConf struct {
	DBName         string
	CollectionName string
	ConvMsgLimit   int
}

type MongoDBSessionService struct {
	client *mongo.Client
	conf   MongoDBSessionServiceConf
}

type mongoSessionDocument struct {
	SessionID string    `bson:"sessionID"`
	Messages  []Message `bson:"messages"`
	CreatedAt time.Time
}

func NewMongoDBSession(client *mongo.Client, conf MongoDBSessionServiceConf) (*MongoDBSessionService, error) {
	return &MongoDBSessionService{
		client: client,
		conf:   conf,
	}, nil
}

func (s *MongoDBSessionService) GetConversationBySessionId(sessionId string) ([]Message, error) {
	collection := s.client.Database(s.conf.DBName).Collection(s.conf.CollectionName)

	var doc mongoSessionDocument
	err := collection.FindOne(context.TODO(), bson.M{"sessionID": sessionId}).Decode(&doc)
	if err != nil {
		return nil, errors.SessionNotFoundError{Message: fmt.Sprintf("sessionID: %s not found", sessionId)}
	}

	// Apply convMsgLimit if set (>0)
	if s.conf.ConvMsgLimit > 0 && len(doc.Messages) > s.conf.ConvMsgLimit {
		start := len(doc.Messages) - s.conf.ConvMsgLimit
		return doc.Messages[start:], nil
	}

	return doc.Messages, nil
}

func (s *MongoDBSessionService) CreateNewSession() (string, error) {
	sessionId := uuid.NewString()
	document := mongoSessionDocument{
		SessionID: sessionId,
		Messages:  make([]Message, 0),
		CreatedAt: time.Now(),
	}

	collection := s.client.Database(s.conf.DBName).Collection(s.conf.CollectionName)
	_, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (s *MongoDBSessionService) AddMessage(sessionId string, msg Message) error {
	collection := s.client.Database(s.conf.DBName).Collection(s.conf.CollectionName)

	update := bson.M{
		"$push": bson.M{"messages": msg},
	}
	opts := options.UpdateOne().SetUpsert(false) // we don't create a new session here

	res, err := collection.UpdateOne(context.TODO(), bson.M{"sessionID": sessionId}, update, opts)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.SessionNotFoundError{Message: fmt.Sprintf("sessionID: %s not found", sessionId)}
	}

	return nil
}
