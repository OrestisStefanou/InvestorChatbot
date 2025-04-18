package services

import (
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
	"time"
)

type StockOverviewDataService interface {
	GetStockProfile(symbol string) (domain.StockProfile, error)
	GetFinancialRatios(symbol string) ([]domain.FinancialRatios, error)
	GetStockForecast(symbol string) (domain.StockForecast, error)
}

type stockOverviewContext struct {
	currentDate          string
	symbol               string
	stockProfile         domain.StockProfile
	stockFinancialRatios []domain.FinancialRatios
	stockForecast        domain.StockForecast
}

type StockOverviewRag struct {
	dataService StockOverviewDataService
	llm         Llm
}

func NewStockOverviewRag(llm Llm, stockOverviewDataService StockOverviewDataService) (*StockOverviewRag, error) {
	return &StockOverviewRag{llm: llm, dataService: stockOverviewDataService}, nil
}

func (rag StockOverviewRag) createRagContext(symbols []string) (string, error) {
	ragContext := make([]stockOverviewContext, 0, len(symbols))

	for _, symbol := range symbols {
		context := stockOverviewContext{}
		context.symbol = symbol
		context.currentDate = time.Now().Format("2006-01-02")

		stockProfile, err := rag.dataService.GetStockProfile(symbol)
		if err != nil {
			return "", &DataServiceError{Message: fmt.Sprintf("GetStockProfile failed: %s", err)}
		}
		context.stockProfile = stockProfile

		stockFinancialRatios, err := rag.dataService.GetFinancialRatios(symbol)
		if err != nil {
			return "", &DataServiceError{Message: fmt.Sprintf("GetFinancialRatios failed: %s", err)}
		}
		context.stockFinancialRatios = stockFinancialRatios

		stockForecast, err := rag.dataService.GetStockForecast(symbol)
		if err != nil {
			return "", &DataServiceError{Message: fmt.Sprintf("GetStockForecast failed: %s", err)}
		}
		context.stockForecast = stockForecast

		ragContext = append(ragContext, context)
	}

	return fmt.Sprintf("%+v\n", ragContext), nil
}

func (rag StockOverviewRag) GenerateRagResponse(conversation []Message, tags Tags, responseChannel chan<- string) error {
	// Format the prompt to contain the neccessary context
	ragContext, err := rag.createRagContext(tags.StockSymbols)
	if err != nil {
		return err
	}

	prompt := fmt.Sprintf(prompts.StockOverviewPrompt, ragContext)
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
