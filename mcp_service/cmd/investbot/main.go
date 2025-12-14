package main

import (
	"context"
	"fmt"
	restHandlers "investbot/pkg/api/rest/handlers"
	"investbot/pkg/config"
	"investbot/pkg/gemini"
	"investbot/pkg/llama"
	"investbot/pkg/marketDataScraper"
	"investbot/pkg/openAI"
	"investbot/pkg/repositories"
	"investbot/pkg/services"
	"log"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	case config.GEMINI:
		llmConfig := gemini.GeminiLlmConfig{
			ModelName:   conf.GeminiModelName,
			Temperature: conf.BaseLlmTemperature,
			ApiKey:      conf.GeminiKey,
		}
		llm, err = gemini.NewGeminiLLM(llmConfig)
	default:
		err = fmt.Errorf("no valid llm provider found")
	}

	return llm, err
}

func initMongoClient(uri string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func main() {
	e := echo.New()

	conf, _ := config.LoadConfig()

	llm, err := getLlm(conf)
	if err != nil {
		log.Fatal(err)
	}

	var (
		userContextRepository  services.UserContextRepository
		topicAndTagsRepository services.TopicAndTagsRepository
		ragResponsesRepository services.RagResponsesRepository
		sessionService         services.SessionService
		mongoClient            *mongo.Client
	)

	// Create Mongo client only once if needed
	if conf.DatabaseProvider == config.MONGO_DB || conf.SessionStorageProvider == config.MONGO_DB_STORAGE {
		mongoClient, err = initMongoClient(conf.MongoDBConf.Uri)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err = mongoClient.Disconnect(context.TODO()); err != nil {
				log.Fatal(err)
			}
		}()
	}

	// User context repository
	switch conf.DatabaseProvider {
	case config.BADGER_DB:
		db, err := badger.Open(badger.DefaultOptions(conf.BadgerDbPath))
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		userContextRepository, err = repositories.NewUserContextRepository(db)
		if err != nil {
			log.Fatal(err)
		}

		topicAndTagsRepository, err = repositories.NewTopicAndTagsBagderRepo(db)
		if err != nil {
			log.Fatal(err)
		}

		ragResponsesRepository, err = repositories.NewRagResponsesBadgerRepo(db)
		if err != nil {
			log.Fatal(err)
		}

	case config.MONGO_DB:
		userContextRepository, err = repositories.NewUserContextMongoRepo(
			mongoClient,
			conf.MongoDBConf.DBName,
			conf.MongoDBConf.UserContextColletionName,
		)
		if err != nil {
			log.Fatal(err)
		}

		topicAndTagsRepository, err = repositories.NewTopicAndTagsMongoRepo(
			mongoClient,
			conf.MongoDBConf.DBName,
			conf.MongoDBConf.TopicAndTagsCollectionName,
		)
		if err != nil {
			log.Fatal(err)
		}

		ragResponsesRepository, err = repositories.NewRagResponsesMongoRepo(
			mongoClient,
			conf.MongoDBConf.DBName,
			conf.MongoDBConf.RagResponsesCollectionName,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Session service
	switch conf.SessionStorageProvider {
	case config.IN_MEMORY_STORAGE:
		sessionService, _ = services.NewInMemorySession(conf.ConvMsgLimit)
	case config.MONGO_DB_STORAGE:
		sessionService, err = services.NewMongoDBSession(
			mongoClient,
			services.MongoDBSessionServiceConf{
				DBName:         conf.MongoDBConf.DBName,
				CollectionName: conf.MongoDBConf.SessionCollectionName,
				ConvMsgLimit:   conf.ConvMsgLimit,
			},
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Setup cache and data services
	cache, _ := services.NewBadgerCacheService()
	dataService := marketDataScraper.NewMarketDataScraperWithCache(cache, conf)
	userContextService, _ := services.NewUserContextService(userContextRepository)

	// Set up rags
	sectorRag, _ := services.NewSectorRag(llm, dataService, userContextService, ragResponsesRepository)
	educationRag, _ := services.NewEducationRag(llm, userContextService, ragResponsesRepository)
	industryRag, _ := services.NewIndustryRag(llm, dataService)
	stockOverviewRag, _ := services.NewStockOverviewRag(llm, dataService, userContextService, ragResponsesRepository)
	stockFinancialsRag, _ := services.NewStockFinancialsRag(llm, dataService, userContextService, ragResponsesRepository)
	etfRag, _ := services.NewEtfRag(llm, dataService, userContextService, ragResponsesRepository)
	newsRag, _ := services.NewMarketNewsRag(llm, dataService, userContextService, ragResponsesRepository)
	followUpQuestionsRag, _ := services.NewFollowUpQuestionsRag(llm, ragResponsesRepository)

	topicToRagMap := map[services.Topic]services.Rag{
		services.SECTORS:          sectorRag,
		services.EDUCATION:        educationRag,
		services.INDUSTRIES:       industryRag,
		services.STOCK_OVERVIEW:   stockOverviewRag,
		services.STOCK_FINANCIALS: stockFinancialsRag,
		services.ETFS:             etfRag,
		services.NEWS:             newsRag,
	}

	// Set up core services
	topicExtractorService, _ := services.NewTopicExtractor(llm, userContextService, ragResponsesRepository)
	tagExtractorService, _ := services.NewTagExtractor(llm, dataService, userContextService, ragResponsesRepository)
	chatService, _ := services.NewChatService(
		topicToRagMap,
		sessionService,
		topicExtractorService,
		tagExtractorService,
		topicAndTagsRepository,
	)
	followUpQuestionsService, _ := services.NewFollowUpQuestionsService(sessionService, followUpQuestionsRag)
	faqService, _ := services.NewFaqService(conf.FaqLimit)
	tickerService, _ := services.NewTickerService(dataService)
	etfService, _ := services.NewEtfService(dataService)
	superInvestorService, _ := services.NewSuperInvestorService(dataService)

	// Set up rest api handlers
	chatHandler, _ := restHandlers.NewChatHandler(chatService)
	sessionHandler, _ := restHandlers.NewSessionHandler(sessionService)
	followUpQuestionsHandler, _ := restHandlers.NewFollowUpQuestionsHandler(followUpQuestionsService, conf.FollowUpQuestionsNum)
	faqHandler, _ := restHandlers.NewFaqHandler(faqService)
	tickerHandler, _ := restHandlers.NewTickerHandler(tickerService)
	etfHandler, _ := restHandlers.NewEtfHandler(etfService)
	superInvestorHandler, _ := restHandlers.NewSuperInvestorHandler(superInvestorService)
	sectorHandler, _ := restHandlers.NewSectorHandler(dataService)
	topicHandler, _ := restHandlers.NewTopicHandler()
	userContextHandler, _ := restHandlers.NewUserContextHandler(userContextService)

	// Set up api routes
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
