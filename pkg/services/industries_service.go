package services

import (
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
)

type IndustryDataService interface {
	GetIndustryStocks(industry string) ([]domain.IndustryStock, error)
	GetIndustries() ([]domain.Industry, error)
}

type industryContext struct {
	industry domain.Industry
	//industryStocks []domain.IndustryStock
}

type IndustryRag struct {
	dataService IndustryDataService
	llm         Llm
}

func NewIndustryRag(llm Llm, industryDataService IndustryDataService) (*IndustryRag, error) {
	return &IndustryRag{llm: llm, dataService: industryDataService}, nil
}

func (rag IndustryRag) createRagContext() (string, error) {
	var ragContext string
	industries, err := rag.dataService.GetIndustries()
	if err != nil {
		return ragContext, &DataServiceError{Message: fmt.Sprintf("GetIndustries failed: %s", err)}
	}

	for i := 0; i < len(industries); i++ {
		industry := industries[i]
		// FOR NOW WE COMMENT OUT THE ONE BELOW TO AVOID MAKING 150+ http calls
		// We can have extra filtering if we want to have context about industry stocks
		// by providing the specific industry
		// industryStocks, err := rag.dataService.GetIndustryStocks(industry.UrlName)
		// if err != nil {
		// 	return ragContext, &DataServiceError{Message: fmt.Sprintf("GetIndustryStocks failed: %s", err)}
		// }

		context := industryContext{
			industry: industry,
			//industryStocks: industryStocks[:2],
		}
		ragContext += fmt.Sprintf("%+v\n", context)
	}

	return ragContext, nil
}

func (rag IndustryRag) GenerateRagResponse(conversation []Message, responseChannel chan<- string) error {
	// Format the prompt to contain the neccessary context
	ragContext, err := rag.createRagContext()
	if err != nil {
		return err
	}
	prompt := fmt.Sprintf(prompts.IndustriesPrompt, ragContext)
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
