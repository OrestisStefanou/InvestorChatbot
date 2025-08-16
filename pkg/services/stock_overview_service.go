package services

import (
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
	"sync"
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
	BaseRag
	dataService        StockOverviewDataService
	userContextService UserContextDataService
}

func NewStockOverviewRag(
	llm Llm,
	stockOverviewDataService StockOverviewDataService,
	userContextService UserContextDataService,
	responsesStore RagResponsesRepository,
) (*StockOverviewRag, error) {
	rag := StockOverviewRag{
		dataService:        stockOverviewDataService,
		userContextService: userContextService,
	}
	rag.llm = llm
	rag.topic = STOCK_OVERVIEW
	rag.responseStore = responsesStore

	return &rag, nil
}

func (rag StockOverviewRag) createRagContext(symbols []string) (string, error) {
	ragContext := make([]stockOverviewContext, 0, len(symbols))

	for _, symbol := range symbols {
		context := stockOverviewContext{
			symbol:      symbol,
			currentDate: time.Now().Format("2006-01-02"),
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		var fetchErr error

		wg.Add(8) // 3 main + 5 historical

		// Fetch stock profile
		go func() {
			defer wg.Done()
			stockProfile, err := rag.dataService.GetStockProfile(symbol)
			if err != nil {
				mu.Lock()
				fetchErr = &DataServiceError{Message: fmt.Sprintf("GetStockProfile failed: %s", err)}
				mu.Unlock()
				return
			}
			mu.Lock()
			context.stockProfile = stockProfile
			mu.Unlock()
		}()

		// Fetch financial ratios
		go func() {
			defer wg.Done()
			stockFinancialRatios, err := rag.dataService.GetFinancialRatios(symbol)
			if err != nil {
				mu.Lock()
				fetchErr = &DataServiceError{Message: fmt.Sprintf("GetFinancialRatios failed: %s", err)}
				mu.Unlock()
				return
			}
			mu.Lock()
			context.stockFinancialRatios = stockFinancialRatios
			mu.Unlock()
		}()

		// Fetch forecast
		go func() {
			defer wg.Done()
			stockForecast, err := rag.dataService.GetStockForecast(symbol)
			if err != nil {
				mu.Lock()
				fetchErr = &DataServiceError{Message: fmt.Sprintf("GetStockForecast failed: %s", err)}
				mu.Unlock()
				return
			}
			mu.Lock()
			context.stockForecast = stockForecast
			mu.Unlock()
		}()

		// Fetch historical performance
		periods := []domain.Period{domain.Period5D, domain.Period1M, domain.Period6M, domain.Period1Y, domain.Period5Y}
		performanceList := make([]stockHistoricalPerformance, 5)

		for i, period := range periods {
			index, perfomancePeriod := i, period // capture loop variables for go routines
			go func() {
				defer wg.Done()
				perf, err := rag.dataService.GetHistoricalPrices(symbol, domain.Stock, perfomancePeriod)
				if err != nil {
					mu.Lock()
					fetchErr = &DataServiceError{Message: fmt.Sprintf("GetHistoricalPrices for %s failed: %s", period, err)}
					mu.Unlock()
					return
				}
				entry := stockHistoricalPerformance{
					period:           perf.Period,
					percentageChange: perf.PercentageChange,
				}
				mu.Lock()
				performanceList[index] = entry
				mu.Unlock()
			}()
		}

		wg.Wait()

		// Check if any fetch failed
		if fetchErr != nil {
			return "", fetchErr
		}

		context.historicalPerformance = performanceList
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

	return rag.GenerateLllmResponse(prompt, conversation, responseChannel)
}
