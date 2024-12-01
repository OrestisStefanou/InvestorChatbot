package services

import "fmt"

type Rag interface {
	GenerateRagResponse(sessionId string, question string, responseChannel chan<- string) error
}

type Topic string

const (
	EDUCATION Topic = "education"
	SECTORS   Topic = "sectors"
)

type ChatService struct {
	topicToRagMap map[Topic]Rag
	database      interface{}
}

func NewChatService(topicToRagMap map[Topic]Rag, database interface{}) (*ChatService, error) {
	return &ChatService{
		topicToRagMap: topicToRagMap,
		database:      database,
	}, nil
}

func (s *ChatService) GenerateResponse(topic Topic, sessionId string, question string, responseChannel chan<- string) error {
	rag, found := s.topicToRagMap[topic]

	if !found {
		return fmt.Errorf("Rag for topic %s not found", topic)
	}

	if err := rag.GenerateRagResponse(sessionId, question, responseChannel); err != nil {
		return err
	}

	return nil
}
