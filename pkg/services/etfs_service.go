package services

import (
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
)

type EtfDataService interface {
	GetEtfs() ([]domain.Etf, error)
	GetEtfOverview(symbol string) (domain.EtfOverview, error)
}

type etfContext struct {
	etf domain.Etf
}

type etfOverviewContext struct {
	etfOverview domain.EtfOverview
}

type EtfRag struct {
	dataService EtfDataService
	llm         Llm
}

func NewEtfRag(llm Llm, etfDataService EtfDataService) (*EtfRag, error) {
	return &EtfRag{llm: llm, dataService: etfDataService}, nil
}

func (rag EtfRag) createRagContext(etfSymbol string) (string, error) {
	var ragContext string

	if etfSymbol != "" {
		etfOverview, err := rag.dataService.GetEtfOverview(etfSymbol)
		if err != nil {
			return ragContext, &DataServiceError{Message: fmt.Sprintf("GetEtfOverview failed: %s", err)}
		}
		context := etfOverviewContext{etfOverview: etfOverview}
		return fmt.Sprintf("%+v\n", context), nil
	}

	etfs, err := rag.dataService.GetEtfs()
	if err != nil {
		return ragContext, &DataServiceError{Message: fmt.Sprintf("GetEtfs failed: %s", err)}
	}

	var context etfContext
	for _, etf := range etfs {
		// Because of the huge number of ETFs we keep only the ones
		// with aum higher than 500M
		if etf.Aum > 500000000 {
			context.etf = etf
			ragContext += fmt.Sprintf("%+v\n", context)
		}
	}

	return ragContext, nil
}

func (rag EtfRag) GenerateRagResponse(conversation []Message, tags Tags, responseChannel chan<- string) error {
	// Format the prompt to contain the neccessary context
	ragContext, err := rag.createRagContext(tags.EtfSymbol)
	if err != nil {
		return err
	}
	prompt := fmt.Sprintf(prompts.EtfsPrompt, ragContext)
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
