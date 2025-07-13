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
	GetHistoricalPrices(ticker string, assetClass domain.AssetClass, period domain.Period) (domain.HistoricalPrices, error)
}

type stockHistoricalPerformance struct {
	period           domain.Period
	percentageChange float64
}

type stockOverviewContext struct {
	currentDate           string
	symbol                string
	stockProfile          domain.StockProfile
	stockFinancialRatios  []domain.FinancialRatios
	stockForecast         domain.StockForecast
	historicalPerformance []stockHistoricalPerformance
}

type StockOverviewRag struct {
	dataService        StockOverviewDataService
	llm                Llm
	userContextService UserContextDataService
}

func NewStockOverviewRag(
	llm Llm,
	stockOverviewDataService StockOverviewDataService,
	userContextService UserContextDataService,
) (*StockOverviewRag, error) {
	return &StockOverviewRag{
		llm:                llm,
		dataService:        stockOverviewDataService,
		userContextService: userContextService,
	}, nil
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

		context.historicalPerformance = make([]stockHistoricalPerformance, 0, 5)
		performance5D, err := rag.dataService.GetHistoricalPrices(symbol, domain.Stock, domain.Period5D)
		if err != nil {
			return "", &DataServiceError{Message: fmt.Sprintf("GetHistoricalPrices for 5D failed: %s", err)}
		}
		context.historicalPerformance = append(
			context.historicalPerformance,
			stockHistoricalPerformance{
				period:           performance5D.Period,
				percentageChange: performance5D.PercentageChange,
			},
		)

		performance1M, err := rag.dataService.GetHistoricalPrices(symbol, domain.Stock, domain.Period1M)
		if err != nil {
			return "", &DataServiceError{Message: fmt.Sprintf("GetHistoricalPrices for 1M failed: %s", err)}
		}
		context.historicalPerformance = append(
			context.historicalPerformance,
			stockHistoricalPerformance{
				period:           performance1M.Period,
				percentageChange: performance1M.PercentageChange,
			},
		)

		performance6M, err := rag.dataService.GetHistoricalPrices(symbol, domain.Stock, domain.Period6M)
		if err != nil {
			return "", &DataServiceError{Message: fmt.Sprintf("GetHistoricalPrices for 6M failed: %s", err)}
		}
		context.historicalPerformance = append(
			context.historicalPerformance,
			stockHistoricalPerformance{
				period:           performance6M.Period,
				percentageChange: performance6M.PercentageChange,
			},
		)

		performance1Y, err := rag.dataService.GetHistoricalPrices(symbol, domain.Stock, domain.Period1Y)
		if err != nil {
			return "", &DataServiceError{Message: fmt.Sprintf("GetHistoricalPrices for 1Y failed: %s", err)}
		}
		context.historicalPerformance = append(
			context.historicalPerformance,
			stockHistoricalPerformance{
				period:           performance1Y.Period,
				percentageChange: performance1Y.PercentageChange,
			},
		)

		performance5Y, err := rag.dataService.GetHistoricalPrices(symbol, domain.Stock, domain.Period5Y)
		if err != nil {
			return "", &DataServiceError{Message: fmt.Sprintf("GetHistoricalPrices for 5y failed: %s", err)}
		}
		context.historicalPerformance = append(
			context.historicalPerformance,
			stockHistoricalPerformance{
				period:           performance5Y.Period,
				percentageChange: performance5Y.PercentageChange,
			},
		)

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

	var userContext domain.UserContext
	if tags.UserID != "" {
		userContext, err = rag.userContextService.GetUserContext(tags.UserID)
		if err != nil {
			return err
		}
	}

	prompt := fmt.Sprintf(prompts.StockOverviewPrompt, ragContext, userContext)
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
