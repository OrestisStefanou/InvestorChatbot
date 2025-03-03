package services

import (
	"fmt"
	"investbot/pkg/errors"
)

type Tags struct {
	SectorName      string
	IndustryName    string
	StockSymbol     string
	BalanceSheet    bool
	IncomeStatement bool
	CashFlow        bool
}

type Rag interface {
	GenerateRagResponse(conversation []Message, tags Tags, responseChannel chan<- string) error
}

type Topic string

const (
	EDUCATION        Topic = "education"
	SECTORS          Topic = "sectors"
	INDUSTRIES       Topic = "industries"
	STOCK_OVERVIEW   Topic = "stock_overview"
	STOCK_FINANCIALS Topic = "stock_financials"
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

func (s *ChatService) GenerateResponse(
	topic Topic,
	tags Tags,
	sessionId string,
	question string,
	responseChannel chan<- string,
) error {
	rag, found := s.topicToRagMap[topic]

	if !found {
		// Use default RAG in this case?
		return fmt.Errorf("Rag for topic %s not found", topic)
	}

	conversation, err := s.sessionService.GetConversationBySessionId(sessionId)
	if err != nil {
		return &errors.SessionNotFoundError{
			Message: fmt.Sprintf("Conversation for session id: %s not found", sessionId),
		}
	}

	questionMessage := Message{
		Role: User, Content: question,
	}
	conversation = append(conversation, questionMessage)

	var responseMessage string
	chunkChannel := make(chan string)
	errorChannel := make(chan error, 1)

	go func() {
		if err := rag.GenerateRagResponse(conversation, tags, chunkChannel); err != nil {
			errorChannel <- err
		}
		close(errorChannel)
	}()

	shouldExit := false
	for !shouldExit {
		select {
		case chunk, isOpen := <-chunkChannel:
			if !isOpen {
				fmt.Printf("FINAL RESPONSE\n %s", responseMessage)
				shouldExit = true
				close(responseChannel)
				continue
			}
			responseMessage += chunk
			responseChannel <- chunk
		case err := <-errorChannel:
			if err != nil {
				return err
			}
		}
	}

	conversation = append(conversation, Message{Role: Assistant, Content: responseMessage})
	s.sessionService.SetConversationForSessionId(conversation, sessionId)

	return nil
}
