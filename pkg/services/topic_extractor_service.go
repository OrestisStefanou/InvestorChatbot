package services

import (
	"fmt"
	"investbot/pkg/services/prompts"
)

type TopicExtractor struct {
	llm Llm
}

func NewTopicExtractor(llm Llm) (*TopicExtractor, error) {
	return &TopicExtractor{llm: llm}, nil
}

func (te TopicExtractor) ExtractTopic(conversation []Message) (Topic, error) {
	prompt := fmt.Sprintf(prompts.TopicExtractorPrompt, conversation)

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

	_, found := topics[Topic(responseMessage)]
	if !found {
		return "", fmt.Errorf("%s is not a valid topic", responseMessage)
	}

	return Topic(responseMessage), nil
}
