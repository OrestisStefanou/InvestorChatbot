package services

import (
	"encoding/json"
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
	"log"
	"strings"
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

type llmTopicResponse struct {
	Topic string `json:"topic"`
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

	go func() {
		storeErr := te.responseStore.StoreRagResponse(
			te.llm.GetLlmName(),
			"ExtractTopic",
			[]Message{promptMsg},
			responseMessage,
		)
		if storeErr != nil {
			log.Printf("Failed to store topic extraction rag response: %s", storeErr.Error())
		}
	}()

	// Validate the response against known topics
	topics := map[Topic]any{
		EDUCATION:        nil,
		SECTORS:          nil,
		STOCK_OVERVIEW:   nil,
		STOCK_FINANCIALS: nil,
		ETFS:             nil,
		NEWS:             nil,
	}

	// Strip formatting artifacts from the response(in case they exist)
	strippedLlmResponse := strings.TrimPrefix(responseMessage, "```json\n")
	strippedLlmResponse = strings.TrimSuffix(strippedLlmResponse, "\n```")

	var topicResponse llmTopicResponse
	err = json.Unmarshal([]byte(strippedLlmResponse), &topicResponse)
	if err != nil {
		return "", err
	}

	if _, found := topics[Topic(topicResponse.Topic)]; !found {
		return "", fmt.Errorf("%s is not a valid topic", topicResponse.Topic)
	}

	return Topic(topicResponse.Topic), nil
}
