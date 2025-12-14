package services

import (
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
)

type EducationRag struct {
	BaseRag
	userContextService UserContextDataService
}

func NewEducationRag(
	llm Llm,
	userContextService UserContextDataService,
	responsesStore RagResponsesRepository,
) (*EducationRag, error) {
	rag := EducationRag{userContextService: userContextService}
	rag.llm = llm
	rag.topic = EDUCATION
	rag.responseStore = responsesStore

	return &rag, nil
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

	return rag.GenerateLllmResponse(prompt, conversation, responseChannel)
}
