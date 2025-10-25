package services

import (
	"encoding/json"
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
	"log"
	"strings"
)

type MarketDataService interface {
	GetSectors() ([]domain.Sector, error)
	GetTickers() ([]domain.Ticker, error)
	GetEtfs() ([]domain.Etf, error)
}

type TagExtractor struct {
	llm                Llm
	marketDataService  MarketDataService
	userContextService UserContextDataService
	responseStore      RagResponsesRepository
}

type llmTagExtractorResponse struct {
	Sector          string   `json:"sector_name"`
	StockSymbols    []string `json:"stock_symbols"`
	BalanceSheet    bool     `json:"balance_sheet"`
	IncomeStatement bool     `json:"income_statement"`
	CashFlow        bool     `json:"cash_flow"`
	EtfSymbols      []string `json:"etf_symbols"`
}

func NewTagExtractor(
	llm Llm,
	marketDataService MarketDataService,
	userContextService UserContextDataService,
	responseStore RagResponsesRepository,
) (*TagExtractor, error) {
	return &TagExtractor{
		llm:                llm,
		marketDataService:  marketDataService,
		userContextService: userContextService,
		responseStore:      responseStore,
	}, nil
}

func (te TagExtractor) ExtractTags(topic Topic, conversation []Message, userID string) (Tags, error) {
	var tags Tags
	var err error
	var userContext domain.UserContext
	if userID != "" {
		userContext, err = te.userContextService.GetUserContext(userID)
		if err != nil {
			return Tags{}, err
		}
	}
	switch topic {
	case SECTORS:
		tags, err = te.extractSectorTags(conversation, userContext)
	case STOCK_OVERVIEW:
		tags, err = te.extractStockOverviewTags(conversation, userContext)
	case STOCK_FINANCIALS:
		tags, err = te.extractStockFinancialsTags(conversation, userContext)
	case ETFS:
		tags, err = te.extractEtfTags(conversation, userContext)
	case NEWS:
		tags, err = te.extractMarketNewsTags(conversation, userContext)
	}
	return tags, err
}

func (te TagExtractor) extractSectorTags(conversation []Message, userContext domain.UserContext) (Tags, error) {
	sectors, err := te.marketDataService.GetSectors()
	if err != nil {
		return Tags{}, err
	}

	var sectorsPlaceholderString string
	for _, s := range sectors {
		sectorsPlaceholderString += fmt.Sprintf("%s\n", s.UrlName)
	}

	prompt := fmt.Sprintf(prompts.SectorTagExtractorPrompt, sectorsPlaceholderString, userContext, conversation)
	llmResponse, err := te.getLlmResponse(prompt)
	if err != nil {
		return Tags{}, err
	}

	var result llmTagExtractorResponse
	err = json.Unmarshal([]byte(llmResponse), &result)
	if err != nil {
		return Tags{}, err
	}

	return Tags{SectorName: result.Sector}, nil
}

func (te TagExtractor) extractStockOverviewTags(conversation []Message, userContext domain.UserContext) (Tags, error) {
	stockSymbols, err := te.marketDataService.GetTickers()
	if err != nil {
		return Tags{}, err
	}

	prompt := fmt.Sprintf(prompts.StockOverviewTagExtractorPrompt, stockSymbols, userContext, conversation)
	llmResponse, err := te.getLlmResponse(prompt)
	if err != nil {
		return Tags{}, err
	}

	var result llmTagExtractorResponse
	err = json.Unmarshal([]byte(llmResponse), &result)
	if err != nil {
		return Tags{}, err
	}

	return Tags{StockSymbols: result.StockSymbols}, nil
}

func (te TagExtractor) extractStockFinancialsTags(conversation []Message, userContext domain.UserContext) (Tags, error) {
	stockSymbols, err := te.marketDataService.GetTickers()
	if err != nil {
		return Tags{}, err
	}

	prompt := fmt.Sprintf(prompts.StockFinancialsTagExtractorPrompt, stockSymbols, userContext, conversation)
	llmResponse, err := te.getLlmResponse(prompt)
	if err != nil {
		return Tags{}, err
	}

	var result llmTagExtractorResponse
	err = json.Unmarshal([]byte(llmResponse), &result)
	if err != nil {
		return Tags{}, err
	}

	return Tags{
		StockSymbols:    result.StockSymbols,
		BalanceSheet:    result.BalanceSheet,
		CashFlow:        result.CashFlow,
		IncomeStatement: result.IncomeStatement,
	}, nil
}

func (te TagExtractor) extractEtfTags(conversation []Message, userContext domain.UserContext) (Tags, error) {
	etfs, err := te.marketDataService.GetEtfs()
	if err != nil {
		return Tags{}, err
	}

	type etfTicker struct {
		etfName   string
		etfSymbol string
	}

	etfSymbols := make([]etfTicker, 0, len(etfs))
	for _, e := range etfs {
		etfSymbols = append(etfSymbols, etfTicker{etfName: e.Name, etfSymbol: e.Symbol})
	}

	prompt := fmt.Sprintf(prompts.EtfTagExtractorPrompt, etfSymbols, userContext, conversation)
	llmResponse, err := te.getLlmResponse(prompt)
	if err != nil {
		return Tags{}, err
	}

	var result llmTagExtractorResponse
	err = json.Unmarshal([]byte(llmResponse), &result)
	if err != nil {
		return Tags{}, err
	}

	return Tags{EtfSymbols: result.EtfSymbols}, nil
}

func (te TagExtractor) extractMarketNewsTags(conversation []Message, userContext domain.UserContext) (Tags, error) {
	stockSymbols, err := te.marketDataService.GetTickers()
	if err != nil {
		return Tags{}, err
	}

	prompt := fmt.Sprintf(prompts.NewsTagExtractorPrompt, stockSymbols, userContext, conversation)
	llmResponse, err := te.getLlmResponse(prompt)
	if err != nil {
		return Tags{}, err
	}

	var result llmTagExtractorResponse
	err = json.Unmarshal([]byte(llmResponse), &result)
	if err != nil {
		return Tags{}, err
	}

	return Tags{StockSymbols: result.StockSymbols}, nil
}

func (te TagExtractor) getLlmResponse(prompt string) (string, error) {
	promptMsg := Message{
		Role:    User,
		Content: prompt,
	}

	responseMessage, err := streamChunks(
		func(chunkChan chan<- string) error {
			return te.llm.GenerateResponse([]Message{promptMsg}, chunkChan)
		},
		nil, // no need to stream out chunks
	)
	if err != nil {
		return "", err
	}

	go func() {
		storeErr := te.responseStore.StoreRagResponse(
			te.llm.GetLlmName(),
			"ExtractTags",
			[]Message{promptMsg},
			responseMessage,
		)
		if storeErr != nil {
			log.Printf("Failed to store tag extraction rag response: %s", storeErr.Error())
		}
	}()

	// Strip formatting artifacts
	stripped := strings.TrimPrefix(responseMessage, "```json\n")
	stripped = strings.TrimSuffix(stripped, "\n```")

	return stripped, nil
}
