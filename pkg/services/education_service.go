package services

import (
	"fmt"
	"investbot/pkg/services/prompts"
)

type EducationServiceRag struct {
	Llm            Llm
	SessionService SessionService
}

func (rag EducationServiceRag) GenerateRagResponse(sessionId string, question string, responseChannel chan<- string) error {
	conversation, err := rag.SessionService.GetConversationBySessionId(sessionId)
	if err != nil {
		return &SessionServiceError{Message: fmt.Sprintf("GetConversationBySessionId failed: %s", err)}
	}

	if len(conversation) == 0 {
		conversation = append(conversation, map[string]string{"role": "system", "content": prompts.EducationPrompt})
	}
	conversation = append(conversation, map[string]string{"role": "user", "content": question})

	if err := rag.Llm.GenerateResponse(conversation, responseChannel); err != nil {
		return &RagError{Message: fmt.Sprintf("GenerateResponse failed: %s", err)}
	}

	return nil
}
