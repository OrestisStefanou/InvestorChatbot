package main

import (
	"fmt"
	"investbot/pkg/config"
	"investbot/pkg/handlers"
	"investbot/pkg/llama"
	"investbot/pkg/marketDataScraper"
	"investbot/pkg/openAI"
	"investbot/pkg/repositories"
	"investbot/pkg/services"
	"log"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/labstack/echo/v4"
)

func getLlm(conf config.Config) (services.Llm, error) {
	var llm services.Llm
	var err error
	switch conf.LlmProvider {
	case config.OPEN_AI:
		openAiClient, _ := openAI.NewOpenAiClient(conf.OpenAiKey, conf.OpenAiBaseUrl)
		llm, err = openAI.NewOpenAiLLM(conf.OpenAiModelName, openAiClient, float64(conf.BaseLlmTemperature))
	case config.OLLAMA:
		llamaClient, _ := llama.NewOllamaClient(conf.OllamaBaseUrl)
		llm, err = llama.NewLlamaLLM(llama.ModelName(conf.OllamaModelName), llamaClient, conf.BaseLlmTemperature)
	default:
		err = fmt.Errorf("No valid llm provider found")
	}

	return llm, err
}

func main() {
	e := echo.New()

	conf, _ := config.LoadConfig()

	db, err := badger.Open(badger.DefaultOptions(conf.BadgerDbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	llm, err := getLlm(conf)
	if err != nil {
		log.Fatal(err)
	}

	userRepository, err := repositories.NewUserRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	cache, _ := services.NewBadgerCacheService()
	dataService := marketDataScraper.NewMarketDataScraperWithCache(cache, conf)
	sectorRag, _ := services.NewSectorRag(llm, dataService)
	educationRag, _ := services.NewEducationRag(llm)
	industryRag, _ := services.NewIndustryRag(llm, dataService)
	stockOverviewRag, _ := services.NewStockOverviewRag(llm, dataService)
	stockFinancialsRag, _ := services.NewStockFinancialsRag(llm, dataService)
	etfRag, _ := services.NewEtfRag(llm, dataService)
	newsRag, _ := services.NewMarketNewsRag(llm, dataService)
	followUpQuestionsRag, _ := services.NewFollowUpQuestionsRag(llm)
	userService, _ := services.NewUserService(userRepository)

	topicToRagMap := map[services.Topic]services.Rag{
		services.SECTORS:          sectorRag,
		services.EDUCATION:        educationRag,
		services.INDUSTRIES:       industryRag,
		services.STOCK_OVERVIEW:   stockOverviewRag,
		services.STOCK_FINANCIALS: stockFinancialsRag,
		services.ETFS:             etfRag,
		services.NEWS:             newsRag,
	}

	sessionService, _ := services.NewInMemorySession(conf.ConvMsgLimit)
	chatService, _ := services.NewChatService(topicToRagMap, sessionService)
	followUpQuestionsService, _ := services.NewFollowUpQuestionsService(sessionService, followUpQuestionsRag)
	faqService, _ := services.NewFaqService(conf.FaqLimit)
	tickerService, _ := services.NewTickerService(dataService)
	etfService, _ := services.NewEtfService(dataService)
	superInvestorService, _ := services.NewSuperInvestorService(dataService)

	chatHandler, _ := handlers.NewChatHandler(chatService)
	sessionHandler, _ := handlers.NewSessionHandler(sessionService)
	followUpQuestionsHandler, _ := handlers.NewFollowUpQuestionsHandler(followUpQuestionsService, conf.FollowUpQuestionsNum)
	faqHandler, _ := handlers.NewFaqHandler(faqService)
	tickerHandler, _ := handlers.NewTickerHandler(tickerService)
	etfHandler, _ := handlers.NewEtfHandler(etfService)
	superInvestorHandler, _ := handlers.NewSuperInvestorHandler(superInvestorService)
	sectorHandler, _ := handlers.NewSectorHandler(dataService)
	topicHandler, _ := handlers.NewTopicHandler()
	userHandler, _ := handlers.NewUserHandler(userService)

	e.POST("/chat", chatHandler.ChatCompletion)
	e.POST("/session", sessionHandler.CreateNewSession)
	e.POST("/follow_up_questions", followUpQuestionsHandler.GenerateFollowUpQuestions)
	e.GET("/faq", faqHandler.GetFaq)
	e.GET("/tickers", tickerHandler.GetTickers)
	e.GET("/etfs", etfHandler.GetEtfs)
	e.GET("/super_investors", superInvestorHandler.GetSuperInvestors)
	e.GET("/super_investors/portfolio/:super_investor", superInvestorHandler.GetSuperInvestorPortfolio)
	e.GET("/sectors", sectorHandler.GetSectors)
	e.GET("/sectors/stocks/:sector", sectorHandler.GetSectorStocks)
	e.GET("/topics", topicHandler.GetTopics)
	e.POST("/user", userHandler.CreateUser)
	e.Logger.Fatal(e.Start(":1323"))
}
