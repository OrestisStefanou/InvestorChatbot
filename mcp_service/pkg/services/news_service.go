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
	BaseRag
	dataService        NewsDataService
	userContextService UserContextDataService
}

type ragNewsContext struct {
	currentDate string
	news        []domain.NewsArticle
}

func NewMarketNewsRag(
	llm Llm,
	newsDataService NewsDataService,
	userContextService UserContextDataService,
	responsesStore RagResponsesRepository,
) (*MarketNewsRag, error) {
	rag := MarketNewsRag{
		dataService:        newsDataService,
		userContextService: userContextService,
	}
	rag.llm = llm
	rag.topic = NEWS
	rag.responseStore = responsesStore

	return &rag, nil
}

func (rag MarketNewsRag) createRagContext(stockSymbols []string) (string, error) {
	ragContext := make([]ragNewsContext, 0, len(stockSymbols))
	var err error
	context := ragNewsContext{}
	// Keep only the last 20 news
	limit := 20

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

	var userContext domain.UserContext
	if tags.UserID != "" {
		userContext, err = rag.userContextService.GetUserContext(tags.UserID)
		if err != nil {
			return err
		}
	}

	prompt := fmt.Sprintf(prompts.NewsPrompt, ragContext, userContext)

	return rag.GenerateLllmResponse(prompt, conversation, responseChannel)
}
