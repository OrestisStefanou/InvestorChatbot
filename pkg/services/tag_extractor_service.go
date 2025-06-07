package services

import (
	"encoding/json"
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
)

type MarketDataService interface {
	GetSectors() ([]domain.Sector, error)
	GetTickers() ([]domain.Ticker, error)
}

type TagExtractor struct {
	llm               Llm
	marketDataService MarketDataService
}

type llmTagExtractorResponse struct {
	Sector       string   `json:"sector_name"`
	StockSymbols []string `json:"stock_symbols"`
}

func NewTagExtractor(llm Llm, marketDataService MarketDataService) (*TagExtractor, error) {
	return &TagExtractor{llm: llm, marketDataService: marketDataService}, nil
}

func (te TagExtractor) ExtractTags(topic Topic, conversation []Message) (Tags, error) {
	switch topic {
	case SECTORS:
		return te.extractSectorTags(conversation)
	case STOCK_OVERVIEW:
		return te.extractStockOverviewTags(conversation)
	}
	return Tags{}, nil
}

func (te TagExtractor) extractSectorTags(conversation []Message) (Tags, error) {
	// Format the sector tag extraction prompt message
	sectors, err := te.marketDataService.GetSectors()
	if err != nil {
		return Tags{}, err
	}

	var sectorsPlaceholderString string
	for _, s := range sectors {
		sectorsPlaceholderString += fmt.Sprintf("%s\n", s.UrlName)
	}

	prompt := fmt.Sprintf(prompts.SectorTagExtractorPrompt, sectorsPlaceholderString, conversation)
	llmResponse, err := te.getLlmResponse(prompt)

	var result llmTagExtractorResponse
	err = json.Unmarshal([]byte(llmResponse), &result)
	if err != nil {
		return Tags{}, err
	}

	return Tags{SectorName: result.Sector}, nil
}

func (te TagExtractor) extractStockOverviewTags(conversation []Message) (Tags, error) {
	// Format the sector tag extraction prompt message
	stockSymbols, err := te.marketDataService.GetTickers()
	if err != nil {
		return Tags{}, err
	}

	prompt := fmt.Sprintf(prompts.StockOverviewTagExtractorPrompt, stockSymbols, conversation)

	llmResponse, err := te.getLlmResponse(prompt)
	var result llmTagExtractorResponse

	err = json.Unmarshal([]byte(llmResponse), &result)
	if err != nil {
		return Tags{}, err
	}

	return Tags{StockSymbols: result.StockSymbols}, nil
}

func (te TagExtractor) getLlmResponse(prompt string) (string, error) {
	var responseMessage string
	chunkChannel := make(chan string)
	errorChannel := make(chan error, 1)

	promptMsg := Message{
		Role:    User,
		Content: prompt,
	}

	go func() {
		if err := te.llm.GenerateResponse([]Message{promptMsg}, chunkChannel); err != nil {
			errorChannel <- err
		}
		close(errorChannel)
	}()

	shouldExit := false
	for !shouldExit {
		select {
		case chunk, isOpen := <-chunkChannel:
			if !isOpen {
				shouldExit = true
				continue
			}
			responseMessage += chunk
		case err := <-errorChannel:
			if err != nil {
				return "", err
			}
		}
	}

	return responseMessage, nil
}
