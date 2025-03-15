package main

import (
	"investbot/pkg/config"
	"investbot/pkg/handlers"
	"investbot/pkg/marketDataScraper"
	"investbot/pkg/openAI"
	"investbot/pkg/services"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	config, _ := config.LoadConfig()
	openAiClient, _ := openAI.NewOpenAiClient(config.OpenAiKey, config.OpenAiBaseUrl)

	openAiLLM, err := openAI.NewOpenAiLLM(openAI.GPT4_MINI, openAiClient, 0.2)
	if err != nil {
		log.Fatal(err)
	}
	// llamaClient, _ := llama.NewOllamaClient(config.OllamaBaseUrl)
	// llamaLLM, _ := llama.NewLlamaLLM("llama3.2", llamaClient, 0.2)
	dataService := marketDataScraper.MarketDataScraper{}
	sectorRag, _ := services.NewSectorRag(openAiLLM, dataService)
	educationRag, _ := services.NewEducationRag(openAiLLM)
	industryRag, _ := services.NewIndustryRag(openAiLLM, dataService)
	stockOverviewRag, _ := services.NewStockOverviewRag(openAiLLM, dataService)
	stockFinancialsRag, _ := services.NewStockFinancialsRag(openAiLLM, dataService)
	etfRag, _ := services.NewEtfRag(openAiLLM, dataService)
	newsRag, _ := services.NewMarketNewsRag(openAiLLM, dataService)

	topicToRagMap := map[services.Topic]services.Rag{
		services.SECTORS:          sectorRag,
		services.EDUCATION:        educationRag,
		services.INDUSTRIES:       industryRag,
		services.STOCK_OVERVIEW:   stockOverviewRag,
		services.STOCK_FINANCIALS: stockFinancialsRag,
		services.ETFS:             etfRag,
		services.NEWS:             newsRag,
	}

	sessionService, _ := services.NewInMemorySession()
	chatService, _ := services.NewChatService(topicToRagMap, sessionService)

	chatHandler, _ := handlers.NewChatHandler(chatService)
	sessionHandler, _ := handlers.NewSessionHandler(sessionService)

	e.POST("/chat", chatHandler.ChatCompletion)
	e.POST("/session", sessionHandler.CreateNewSession)
	e.Logger.Fatal(e.Start(":1323"))
}
