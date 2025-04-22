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
	AvgPe                float64
	AvgPb                float64
	AvgPfcf              float64
	AvgEvEbitda          float64
	AvgEvEbit            float64
	AvgEvFcf             float64
	AvgDebtEquity        float64
	AvgDebtEbitda        float64
	AvgDebtFcf           float64
	AvgAssetTurnover     float64
	AvgInventoryTurnover float64
	AvgQuickRatio        float64
	AvgCurrentRatio      float64
	AvgRoe               float64
	AvgRoa               float64
	AvgRoic              float64
	AvgEarningYield      float64
	AvgFcfYield          float64
	AvgDividendYield     float64
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

// TODO: Should this be a method on the dataService?
func (rag StockOverviewRag) getStockCompetitors(stockProfile domain.StockProfile, symbol string) ([]domain.IndustryStock, error) {
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

func (rag StockOverviewRag) getCompetitorsAvgFinancialRatios(competitors []domain.IndustryStock) (competitorsAvgFinancialRatios, error) {
	totalRatios := domain.FinancialRatios{}

	for _, competitor := range competitors {
		financialRatios, err := rag.dataService.GetFinancialRatios(competitor.Symbol)
		if err != nil {
			return competitorsAvgFinancialRatios{}, err
		}

		totalRatios.Pe += financialRatios[0].Pe
		totalRatios.Pb += financialRatios[0].Pb
		totalRatios.Pfcf += financialRatios[0].Pfcf
		totalRatios.EvEbitda += financialRatios[0].EvEbitda
		totalRatios.EvEbit += financialRatios[0].EvEbit
		totalRatios.EvFcf += financialRatios[0].EvFcf
		totalRatios.DebtEquity += financialRatios[0].DebtEquity
		totalRatios.DebtEbitda += financialRatios[0].DebtEbitda
		totalRatios.DebtFcf += financialRatios[0].DebtFcf
		totalRatios.AssetTurnover += financialRatios[0].AssetTurnover
		totalRatios.InventoryTurnover += financialRatios[0].InventoryTurnover
		totalRatios.QuickRatio += financialRatios[0].QuickRatio
		totalRatios.CurrentRatio += financialRatios[0].CurrentRatio
		totalRatios.Roe += financialRatios[0].Roe
		totalRatios.Roa += financialRatios[0].Roa
		totalRatios.Roic += financialRatios[0].Roic
		totalRatios.EarningsYield += financialRatios[0].EarningsYield
		totalRatios.FcfYield += financialRatios[0].FcfYield
		totalRatios.DividendYield += financialRatios[0].DividendYield
	}

	avgRatios := competitorsAvgFinancialRatios{}
	avgRatios.AvgPe = totalRatios.Pe / float64(len(competitors))
	avgRatios.AvgPb = totalRatios.Pb / float64(len(competitors))
	avgRatios.AvgPfcf = totalRatios.Pfcf / float64(len(competitors))
	avgRatios.AvgEvEbitda = totalRatios.EvEbitda / float64(len(competitors))
	avgRatios.AvgEvEbit = totalRatios.EvEbit / float64(len(competitors))
	avgRatios.AvgEvFcf = totalRatios.EvFcf / float64(len(competitors))
	avgRatios.AvgDebtEquity = totalRatios.DebtEquity / float64(len(competitors))
	avgRatios.AvgDebtEbitda = totalRatios.DebtEbitda / float64(len(competitors))
	avgRatios.AvgDebtFcf = totalRatios.DebtFcf / float64(len(competitors))
	avgRatios.AvgAssetTurnover = totalRatios.AssetTurnover / float64(len(competitors))
	avgRatios.AvgInventoryTurnover = totalRatios.InventoryTurnover / float64(len(competitors))
	avgRatios.AvgQuickRatio = totalRatios.QuickRatio / float64(len(competitors))
	avgRatios.AvgCurrentRatio = totalRatios.CurrentRatio / float64(len(competitors))
	avgRatios.AvgRoe = totalRatios.Roe / float64(len(competitors))
	avgRatios.AvgRoa = totalRatios.Roa / float64(len(competitors))
	avgRatios.AvgRoic = totalRatios.Roic / float64(len(competitors))
	avgRatios.AvgEarningYield = totalRatios.EarningsYield / float64(len(competitors))
	avgRatios.AvgFcfYield = totalRatios.FcfYield / float64(len(competitors))
	avgRatios.AvgDividendYield = totalRatios.DividendYield / float64(len(competitors))

	return avgRatios, nil
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

		competitors, err := rag.getStockCompetitors(stockProfile, symbol)
		if err != nil {
			fmt.Printf("getStockCompetitors failed: %s", err)
			competitors = []domain.IndustryStock{}
		}
		avgFinancialRatios, err := rag.getCompetitorsAvgFinancialRatios(competitors)
		if err != nil {
			fmt.Printf("getCompetitorsAvgFinancialRatios failed: %s", err)
			avgFinancialRatios = competitorsAvgFinancialRatios{}
		}
		context.competitorsAvgFinancialRatios = avgFinancialRatios

		// TODO: REMOVE THIS
		fmt.Printf("\n\ncontext: %+v\n\n", context)

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
