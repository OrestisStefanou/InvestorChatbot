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

	userContextRepository, err := repositories.NewUserContextRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	cache, _ := services.NewBadgerCacheService()
	dataService := marketDataScraper.NewMarketDataScraperWithCache(cache, conf)
	userContextService, _ := services.NewUserContextService(userContextRepository)
	sectorRag, _ := services.NewSectorRag(llm, dataService, userContextService)
	educationRag, _ := services.NewEducationRag(llm, userContextService)
	industryRag, _ := services.NewIndustryRag(llm, dataService)
	stockOverviewRag, _ := services.NewStockOverviewRag(llm, dataService, userContextService)
	stockFinancialsRag, _ := services.NewStockFinancialsRag(llm, dataService, userContextService)
	etfRag, _ := services.NewEtfRag(llm, dataService, userContextService)
	newsRag, _ := services.NewMarketNewsRag(llm, dataService)
	followUpQuestionsRag, _ := services.NewFollowUpQuestionsRag(llm)

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
	topicExtractorService, _ := services.NewTopicExtractor(llm, userContextService)
	tagExtractorService, _ := services.NewTagExtractor(llm, dataService, userContextService)
	chatService, _ := services.NewChatService(topicToRagMap, sessionService, topicExtractorService, tagExtractorService)
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
	userContextHandler, _ := handlers.NewUserContextHandler(userContextService)

	e.POST("/chat", chatHandler.ChatCompletion)
	e.POST("/chat/extract_topic_and_tags", chatHandler.ExtractTopicAndTags)
	e.POST("/session", sessionHandler.CreateNewSession)
	e.GET("/session/:session_id", sessionHandler.GetSession)
	e.POST("/follow_up_questions", followUpQuestionsHandler.GenerateFollowUpQuestions)
	e.GET("/faq", faqHandler.GetFaq)
	e.GET("/tickers", tickerHandler.GetTickers)
	e.GET("/etfs", etfHandler.GetEtfs)
	e.GET("/super_investors", superInvestorHandler.GetSuperInvestors)
	e.GET("/super_investors/portfolio/:super_investor", superInvestorHandler.GetSuperInvestorPortfolio)
	e.GET("/sectors", sectorHandler.GetSectors)
	e.GET("/sectors/stocks/:sector", sectorHandler.GetSectorStocks)
	e.GET("/topics", topicHandler.GetTopics)
	e.POST("/user_context", userContextHandler.CreateUserContext)
	e.PUT("/user_context", userContextHandler.UpdateUserContext)
	e.GET("/user_context/:user_id", userContextHandler.GetUserContext)
	e.Logger.Fatal(e.Start(":1323"))
}
