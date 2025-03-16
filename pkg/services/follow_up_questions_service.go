package services

import (
	"fmt"
	"investbot/pkg/errors"
	"investbot/pkg/services/prompts"
	"strings"
)

type FollowUpQuestionsRag interface {
	GenerateFollowUpQuestions(conversation []Message, followUpQuestionsNum int) ([]string, error)
}

type FollowUpQuestionsRagImpl struct {
	llm Llm
}

func NewFollowUpQuestionsRag(llm Llm) (*FollowUpQuestionsRagImpl, error) {
	return &FollowUpQuestionsRagImpl{llm: llm}, nil
}

func (rag FollowUpQuestionsRagImpl) GenerateFollowUpQuestions(conversation []Message, followUpQuestionsNum int) ([]string, error) {
	followUpQuestions := make([]string, 0, followUpQuestionsNum)

	prompt := fmt.Sprintf(prompts.FollowUpQuestionsPrompt, followUpQuestionsNum, conversation)
	prompt_msg := Message{
		Role:    System,
		Content: prompt,
	}

	// Add the prompt as the first message in the existing conversation
	conversation_with_prompt := append([]Message{prompt_msg}, conversation...)

	var responseMessage string
	chunkChannel := make(chan string)
	errorChannel := make(chan error, 1)

	go func() {
		if err := rag.llm.GenerateResponse(conversation_with_prompt, chunkChannel); err != nil {
			errorChannel <- err
		}
		close(errorChannel)
	}()

	shouldExit := false
	for !shouldExit {
		select {
		case chunk, isOpen := <-chunkChannel:
			if !isOpen {
				//fmt.Printf("FINAL RESPONSE\n %s", responseMessage)
				shouldExit = true
				continue
			}
			responseMessage += chunk
		case err := <-errorChannel:
			if err != nil {
				return followUpQuestions, err
			}
		}
	}

	// lines := strings.Split(responseMessage, "\n")

	// fmt.Println("\n-----------------------------")

	// for _, line := range lines {
	// 	fmt.Println(line)
	// }

	return strings.Split(responseMessage, "\n"), nil
}

type FollowUpQuestionsService struct {
	sessionService SessionService
	rag            FollowUpQuestionsRag
}

func NewFollowUpQuestionsService(sessionService SessionService, followUpQuestionsRag FollowUpQuestionsRag) (*FollowUpQuestionsService, error) {
	return &FollowUpQuestionsService{sessionService: sessionService, rag: followUpQuestionsRag}, nil
}

func (s FollowUpQuestionsService) GenerateFollowUpQuestions(sessionId string, followUpQuestionsNum int) ([]string, error) {
	conversation, err := s.sessionService.GetConversationBySessionId(sessionId)
	if err != nil {
		return []string{}, &errors.SessionNotFoundError{
			Message: fmt.Sprintf("Conversation for session id: %s not found", sessionId),
		}
	}

	return s.rag.GenerateFollowUpQuestions(conversation, followUpQuestionsNum)
}
