package services

import (
	"fmt"
	"investbot/pkg/errors"
)

type Tags struct {
	SectorName      string
	IndustryName    string
	StockSymbols    []string
	BalanceSheet    bool
	IncomeStatement bool
	CashFlow        bool
	EtfSymbol       string
	UserID          string
}

type Rag interface {
	GenerateRagResponse(conversation []Message, tags Tags, responseChannel chan<- string) error
}

type TopicExtractorService interface {
	ExtractTopic(conversation []Message) (Topic, error)
}

type TagExtractorService interface {
	ExtractTags(topic Topic, conversation []Message) (Tags, error)
}

type Topic string

const (
	EDUCATION        Topic = "education"
	SECTORS          Topic = "sectors"
	INDUSTRIES       Topic = "industries"
	STOCK_OVERVIEW   Topic = "stock_overview"
	STOCK_FINANCIALS Topic = "stock_financials"
	ETFS             Topic = "etfs"
	NEWS             Topic = "news"
)

type ChatService struct {
	topicToRagMap         map[Topic]Rag
	sessionService        SessionService
	topicExtractorService TopicExtractorService
	tagExtractorService   TagExtractorService
}

func NewChatService(
	topicToRagMap map[Topic]Rag,
	sessionService SessionService,
	topicExtractorService TopicExtractorService,
	tagExtractorService TagExtractorService,
) (*ChatService, error) {
	return &ChatService{
		topicToRagMap:         topicToRagMap,
		sessionService:        sessionService,
		topicExtractorService: topicExtractorService,
		tagExtractorService:   tagExtractorService,
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
		return &errors.InvalidTopicError{Message: fmt.Sprintf("Invalid topic %s", topic)}
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
	s.sessionService.AddMessage(sessionId, questionMessage)
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
				fmt.Printf("\n\nFINAL RESPONSE\n %s", responseMessage)
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

	s.sessionService.AddMessage(sessionId, Message{Role: Assistant, Content: responseMessage})

	return nil
}

func (s *ChatService) ExtractTopicAndTags(question string, sessionId string) (Topic, Tags, error) {
	conversation, err := s.sessionService.GetConversationBySessionId(sessionId)
	if err != nil {
		return "", Tags{}, &errors.SessionNotFoundError{
			Message: fmt.Sprintf("Conversation for session id: %s not found", sessionId),
		}
	}

	questionMessage := Message{
		Role: User, Content: question,
	}
	conversation = append(conversation, questionMessage)

	topic, err := s.topicExtractorService.ExtractTopic(conversation)
	if err != nil {
		return "", Tags{}, err
	}

	tags, err := s.tagExtractorService.ExtractTags(topic, conversation)
	if err != nil {
		return "", Tags{}, err
	}

	return topic, tags, nil
}
