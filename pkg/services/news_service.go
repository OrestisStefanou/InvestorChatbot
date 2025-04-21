package services

import (
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
	"time"
)

type NewsDataService interface {
	GetMarketNews() ([]domain.NewsArticle, error)
	GetStockNews(symbol string) ([]domain.NewsArticle, error)
}

type MarketNewsRag struct {
	dataService NewsDataService
	llm         Llm
}

type ragNewsContext struct {
	currentDate string
	news        []domain.NewsArticle
}

func NewMarketNewsRag(llm Llm, newsDataService NewsDataService) (*MarketNewsRag, error) {
	return &MarketNewsRag{llm: llm, dataService: newsDataService}, nil
}

func (rag MarketNewsRag) createRagContext(stockSymbols []string) (string, error) {
	ragContext := make([]ragNewsContext, 0, len(stockSymbols))
	var err error
	context := ragNewsContext{}
	// Keep only the last 10 news
	limit := 10

	if len(stockSymbols) > 0 {
		for _, symbol := range stockSymbols {
			var news []domain.NewsArticle

			news, err = rag.dataService.GetStockNews(symbol)
			if err != nil {
				return "", &DataServiceError{Message: fmt.Sprintf("GetStockNews failed: %s", err)}
			}

			if len(news) < limit {
				limit = len(news)
			}

			context.currentDate = time.Now().Format("2006-01-02")
			context.news = news[:limit]

			ragContext = append(ragContext, context)
		}
	} else {
		news, err := rag.dataService.GetMarketNews()
		if err != nil {
			return "", &DataServiceError{Message: fmt.Sprintf("GetMarketNews failed: %s", err)}
		}
		context.currentDate = time.Now().Format("2006-01-02")
		context.news = news[:limit]

		ragContext = append(ragContext, context)
	}

	return fmt.Sprintf("%+v\n", ragContext), nil
}

func (rag MarketNewsRag) GenerateRagResponse(conversation []Message, tags Tags, responseChannel chan<- string) error {
	// Format the prompt to contain the neccessary context
	ragContext, err := rag.createRagContext(tags.StockSymbols)
	if err != nil {
		return err
	}
	prompt := fmt.Sprintf(prompts.NewsPrompt, ragContext)
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
