package services

import (
	"fmt"
	"investbot/pkg/errors"
	"investbot/pkg/services/prompts"
	"strings"

	"github.com/labstack/gommon/log"
)

type FollowUpQuestionsRag interface {
	GenerateFollowUpQuestions(conversation []Message, followUpQuestionsNum int) ([]string, error)
}

type FollowUpQuestionsRagImpl struct {
	llm           Llm
	responseStore RagResponsesRepository
}

func NewFollowUpQuestionsRag(llm Llm, responsesStore RagResponsesRepository) (*FollowUpQuestionsRagImpl, error) {
	return &FollowUpQuestionsRagImpl{llm: llm, responseStore: responsesStore}, nil
}

func (rag FollowUpQuestionsRagImpl) GenerateFollowUpQuestions(
	conversation []Message,
	followUpQuestionsNum int,
) ([]string, error) {
	followUpQuestions := make([]string, 0, followUpQuestionsNum)

	prompt := fmt.Sprintf(prompts.FollowUpQuestionsPrompt, followUpQuestionsNum, conversation)
	promptMsg := Message{
		Role:    System,
		Content: prompt,
	}

	// Add the prompt as the first message in the existing conversation
	conversationWithPrompt := append([]Message{promptMsg}, conversation...)

	responseMessage, err := streamChunks(
		func(chunkChan chan<- string) error {
			return rag.llm.GenerateResponse(conversationWithPrompt, chunkChan)
		},
		nil, // no need to stream follow-up questions
	)
	if err != nil {
		return followUpQuestions, err
	}

	go func() {
		storeErr := rag.responseStore.StoreRagResponse(
			"FollowUpQuestions",
			conversationWithPrompt,
			responseMessage,
		)
		if storeErr != nil {
			log.Errorf("Failed to store follow up questions rag response: %s", storeErr.Error())
		}
	}()

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
