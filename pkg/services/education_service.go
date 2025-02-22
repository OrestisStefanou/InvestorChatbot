package services

import (
	"fmt"
	"investbot/pkg/services/prompts"
)

type EducationRag struct {
	llm Llm
}

func NewEducationRag(llm Llm) (*EducationRag, error) {
	return &EducationRag{llm: llm}, nil
}

func (rag EducationRag) GenerateRagResponse(conversation []Message, tags Tags, responseChannel chan<- string) error {
	// Add the prompt as the first message in the existing conversation
	prompt_msg := Message{
		Role:    System,
		Content: prompts.EducationPrompt,
	}

	conversation_with_prompt := append([]Message{prompt_msg}, conversation...)

	if err := rag.llm.GenerateResponse(conversation_with_prompt, responseChannel); err != nil {
		return &RagError{Message: fmt.Sprintf("GenerateResponse failed: %s", err)}
	}

	return nil
}
