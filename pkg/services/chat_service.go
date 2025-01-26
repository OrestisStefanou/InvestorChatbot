package services

import (
	"fmt"
	"investbot/pkg/errors"
)

type Rag interface {
	GenerateRagResponse(conversation []Message, responseChannel chan<- string) error
}

type Topic string

const (
	EDUCATION Topic = "education"
	SECTORS   Topic = "sectors"
)

type ChatService struct {
	topicToRagMap  map[Topic]Rag
	sessionService SessionService
}

func NewChatService(topicToRagMap map[Topic]Rag, sessionService SessionService) (*ChatService, error) {
	return &ChatService{
		topicToRagMap:  topicToRagMap,
		sessionService: sessionService,
	}, nil
}

func (s *ChatService) GenerateResponse(topic Topic, sessionId string, question string, responseChannel chan<- string) error {
	rag, found := s.topicToRagMap[topic]

	if !found {
		return fmt.Errorf("Rag for topic %s not found", topic)
	}

	conversation, err := s.sessionService.GetConversationBySessionId(sessionId)

	if err != nil {
		return &errors.SessionNotFoundError{
			Message: fmt.Sprintf("Conversation for session id: %s not found", sessionId),
		}
	}

	if err := rag.GenerateRagResponse(conversation, responseChannel); err != nil {
		return err
	}

	return nil
}
