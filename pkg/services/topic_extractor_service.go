package services

import (
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
	"strings"
)

type TopicExtractor struct {
	llm                Llm
	userContextService UserContextDataService
}

func NewTopicExtractor(llm Llm, userContextService UserContextDataService) (*TopicExtractor, error) {
	return &TopicExtractor{llm: llm, userContextService: userContextService}, nil
}

func (te TopicExtractor) ExtractTopic(conversation []Message, userID string) (Topic, error) {
	var userContext domain.UserContext
	var err error
	if userID != "" {
		userContext, err = te.userContextService.GetUserContext(userID)
		if err != nil {
			return "", err
		}
	}
	prompt := fmt.Sprintf(prompts.TopicExtractorPrompt, userContext, conversation)

	promptMsg := Message{
		Role:    User,
		Content: prompt,
	}

	var responseMessage string
	chunkChannel := make(chan string)
	errorChannel := make(chan error, 1)

	go func() {
		if err := te.llm.GenerateResponse([]Message{promptMsg}, chunkChannel); err != nil {
			errorChannel <- err
		}
		close(errorChannel)
	}()

	shouldExit := false
	for !shouldExit {
		select {
		case chunk, isOpen := <-chunkChannel:
			if !isOpen {
				shouldExit = true
				continue
			}
			responseMessage += chunk
		case err := <-errorChannel:
			if err != nil {
				return "", err
			}
		}
	}
	topics := map[Topic]any{EDUCATION: nil, SECTORS: nil, STOCK_OVERVIEW: nil, STOCK_FINANCIALS: nil, ETFS: nil, NEWS: nil}

	cleanedResponse := strings.ReplaceAll(responseMessage, "\n", "")

	_, found := topics[Topic(cleanedResponse)]
	if !found {
		return "", fmt.Errorf("%s is not a valid topic", cleanedResponse)
	}

	return Topic(cleanedResponse), nil
}
