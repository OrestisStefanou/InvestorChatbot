package services

import (
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
)

type EducationRag struct {
	llm                Llm
	userContextService UserContextDataService
}

func NewEducationRag(llm Llm, userContextService UserContextDataService) (*EducationRag, error) {
	return &EducationRag{llm: llm, userContextService: userContextService}, nil
}

func (rag EducationRag) GenerateRagResponse(conversation []Message, tags Tags, responseChannel chan<- string) error {
	var prompt string
	var userContext domain.UserContext
	var err error

	if tags.UserID != "" {
		userContext, err = rag.userContextService.GetUserContext(tags.UserID)
		if err != nil {
			return err
		}
	}

	prompt = fmt.Sprintf(prompts.EducationPrompt, userContext)
	prompt_msg := Message{
		Role:    System,
		Content: prompt,
	}

	// Add the prompt as the first message in the existing conversation
	conversation_with_prompt := append([]Message{prompt_msg}, conversation...)

	if err := rag.llm.GenerateResponse(conversation_with_prompt, responseChannel); err != nil {
		return &RagError{Message: fmt.Sprintf("GenerateResponse failed: %s", err)}
	}

	return nil
}
