package services

import (
	"encoding/json"
	"fmt"
	"investbot/pkg/errors"
	"investbot/pkg/services/prompts"
	"log"
	"strings"
)

type FollowUpQuestionsRag interface {
	GenerateFollowUpQuestions(conversation []Message, followUpQuestionsNum int) ([]string, error)
}

type FollowUpQuestionsRagImpl struct {
	llm           Llm
	responseStore RagResponsesRepository
}

type llmFollowUpQuestionsResponse struct {
	FollowUpQuestions []string `json:"follow_up_questions"`
}

func NewFollowUpQuestionsRag(llm Llm, responsesStore RagResponsesRepository) (*FollowUpQuestionsRagImpl, error) {
	return &FollowUpQuestionsRagImpl{llm: llm, responseStore: responsesStore}, nil
}

func (rag FollowUpQuestionsRagImpl) GenerateFollowUpQuestions(
	conversation []Message,
	followUpQuestionsNum int,
) ([]string, error) {
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
		return nil, err
	}

	go func() {
		storeErr := rag.responseStore.StoreRagResponse(
			rag.llm.GetLlmName(),
			"FollowUpQuestions",
			conversationWithPrompt,
			responseMessage,
		)
		if storeErr != nil {
			log.Printf("Failed to store follow up questions rag response: %s", storeErr.Error())
		}
	}()

	// Strip formatting artifacts from the response(in case they exist)
	strippedLlmResponse := strings.TrimPrefix(responseMessage, "```json\n")
	strippedLlmResponse = strings.TrimSuffix(strippedLlmResponse, "\n```")

	var followUpsResponse llmFollowUpQuestionsResponse
	err = json.Unmarshal([]byte(strippedLlmResponse), &followUpsResponse)
	if err != nil {
		return nil, err
	}

	return followUpsResponse.FollowUpQuestions, nil
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
