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
	GetIndustryStocks(industry string) ([]domain.IndustryStock, error)
}

type competitorsAvgFinancialRatios struct {
}

type stockOverviewContext struct {
	currentDate                   string
	symbol                        string
	stockProfile                  domain.StockProfile
	stockFinancialRatios          []domain.FinancialRatios
	stockForecast                 domain.StockForecast
	competitorsAvgFinancialRatios competitorsAvgFinancialRatios
}

type StockOverviewRag struct {
	dataService StockOverviewDataService
	llm         Llm
}

func NewStockOverviewRag(llm Llm, stockOverviewDataService StockOverviewDataService) (*StockOverviewRag, error) {
	return &StockOverviewRag{llm: llm, dataService: stockOverviewDataService}, nil
}

func (rag StockOverviewRag) GetStockCompetitors(stockProfile domain.StockProfile, symbol string) ([]domain.IndustryStock, error) {
	// Fetch the industry stocks
	industryStocks, err := rag.dataService.GetIndustryStocks(stockProfile.IndustryUrlName)
	if err != nil {
		return nil, err
	}

	competitors := make([]domain.IndustryStock, 0, 10)
	for i, stock := range industryStocks {
		if stock.Symbol == symbol {
			beforeCount := i
			afterCount := len(industryStocks) - i - 1

			start := 0
			end := 0

			if beforeCount >= 5 && afterCount >= 5 {
				// Ideal: 5 before, 5 after
				start = i - 5
				end = i + 6
			} else if beforeCount < 5 {
				// Not enough before: take up to 10 after
				start = i + 1
				end = min(i+11, len(industryStocks))
			} else if afterCount < 5 {
				// Not enough after: take up to 10 before
				start = max(0, i-10)
				end = i
			}

			// Append from start to end (excluding the matched item)
			competitors = append(competitors, industryStocks[start:end]...)
			break
		}
	}

	return competitors, nil
}

func (rag StockOverviewRag) getCompetitorsAvgFinancialRatios(symbols []string) (competitorsAvgFinancialRatios, error) {
	return competitorsAvgFinancialRatios{}, nil
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
