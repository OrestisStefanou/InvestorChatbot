package services

import (
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
	"strings"

	"github.com/labstack/gommon/log"
)

type TopicExtractor struct {
	llm                Llm
	userContextService UserContextDataService
	responseStore      RagResponsesRepository
}

func NewTopicExtractor(
	llm Llm,
	userContextService UserContextDataService,
	responsesStore RagResponsesRepository,
) (*TopicExtractor, error) {
	return &TopicExtractor{
		llm:                llm,
		userContextService: userContextService,
		responseStore:      responsesStore,
	}, nil
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

	responseMessage, err := streamChunks(
		func(chunkChan chan<- string) error {
			return te.llm.GenerateResponse([]Message{promptMsg}, chunkChan)
		},
		nil, // we don't need to stream this response
	)
	if err != nil {
		return "", err
	}

	// Validate the response against known topics
	topics := map[Topic]any{
		EDUCATION:        nil,
		SECTORS:          nil,
		STOCK_OVERVIEW:   nil,
		STOCK_FINANCIALS: nil,
		ETFS:             nil,
		NEWS:             nil,
	}

	cleanedResponse := strings.ReplaceAll(responseMessage, "\n", "")
	if _, found := topics[Topic(cleanedResponse)]; !found {
		return "", fmt.Errorf("%s is not a valid topic", cleanedResponse)
	}

	go func() {
		storeErr := te.responseStore.StoreRagResponse(
			"ExtractTopic",
			[]Message{promptMsg},
			responseMessage,
		)
		if storeErr != nil {
			log.Errorf("Failed to store topic extraction rag response: %s", storeErr.Error())
		}
	}()

	return Topic(cleanedResponse), nil
}
