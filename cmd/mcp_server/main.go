package main

import (
	"context"
	alphavantage "investbot/pkg/alpha_vantage"
	"investbot/pkg/api/mcp/tools"
	"investbot/pkg/config"
	"investbot/pkg/marketDataScraper"
	"investbot/pkg/repositories"
	"investbot/pkg/services"
	"log"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

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
	conf, _ := config.LoadConfig()

	// Initialize components
	logger := log.New(os.Stdout, "[MCP] ", log.LstdFlags)

	// Create middleware
	loggingMW := NewLoggingMiddleware(logger)

	mcpServer := server.NewMCPServer(
		"Investbot MCP Server", // TODO: Rename this
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(false, true),
		server.WithPromptCapabilities(true),
		server.WithToolHandlerMiddleware(loggingMW.ToolMiddleware),
	)

	var (
		userContextRepository services.UserContextRepository
		mongoClient           *mongo.Client
		err                   error
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

	// Setup cache and data services
	cache, _ := services.NewBadgerCacheService()
	dataService := marketDataScraper.NewMarketDataScraperWithCache(cache, conf)
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

	case config.MONGO_DB:
		userContextRepository, err = repositories.NewUserContextMongoRepo(
			mongoClient,
			conf.MongoDBConf.DBName,
			conf.MongoDBConf.UserContextColletionName,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	alphaVantageClient, _ := alphavantage.NewAlphaVantageClient(conf.AlphaVantageApiKey)

	// Set up services
	tickerService, _ := services.NewTickerService(dataService)
	etfService, _ := services.NewEtfService(dataService)
	superInvestorService, _ := services.NewSuperInvestorService(dataService)
	userContextService, _ := services.NewUserContextService(userContextRepository)

	// Setup tools
	searchStocksTool, _ := tools.NewStockSearchTool(tickerService)
	searchEtfsTool, _ := tools.NewSearchEtfTool(etfService)
	getEtfTool, _ := tools.NewGetEtfTool(etfService)
	getSuperInvestorsTool, _ := tools.NewGetSuperInvestorsTool(superInvestorService)
	getSuperInvestorPortfolioTool, _ := tools.NewGetSuperInvestorPortfolioTool(superInvestorService)
	getMarketNewsTool, _ := tools.NewGetMarketNewsTool(dataService)
	getSectorsTool, _ := tools.NewGetSectorsTool(dataService)
	getSectorStocksTool, _ := tools.NewGetSectorStocksTool(dataService)
	getStockOverviewTool, _ := tools.NewGetStockOverviewTool(dataService)
	getStockFinancialsTool, _ := tools.NewGetStockFinancialsTool(dataService)
	getUserContextTool, _ := tools.NewGetUserContextTool(userContextService)
	updateUserContextTool, _ := tools.NewUpdateUserContextTool(userContextService)
	getEconomicIndicatorTimeSeriesTool, _ := tools.NewGetEconomicIndicatorTimeSeriesTool(alphaVantageClient)

	// Add tools
	mcpServer.AddTool(
		searchStocksTool.GetTool(),
		mcp.NewStructuredToolHandler(searchStocksTool.HandleSearchStocks),
	)

	mcpServer.AddTool(
		searchEtfsTool.GetTool(),
		mcp.NewStructuredToolHandler(searchEtfsTool.HandleSearchEtfs),
	)

	mcpServer.AddTool(
		getEtfTool.GetTool(),
		mcp.NewStructuredToolHandler(getEtfTool.HandleGetEtf),
	)

	mcpServer.AddTool(
		getSuperInvestorsTool.GetTool(),
		mcp.NewStructuredToolHandler(getSuperInvestorsTool.HandleGetSuperInvestors),
	)

	mcpServer.AddTool(
		getSuperInvestorPortfolioTool.GetTool(),
		mcp.NewStructuredToolHandler(getSuperInvestorPortfolioTool.HandleGetSuperInvestorPortfolio),
	)

	mcpServer.AddTool(
		getMarketNewsTool.GetTool(),
		mcp.NewStructuredToolHandler(getMarketNewsTool.HandleGetNews),
	)

	mcpServer.AddTool(
		getSectorsTool.GetTool(),
		mcp.NewStructuredToolHandler(getSectorsTool.HandleGetSectors),
	)

	mcpServer.AddTool(
		getSectorStocksTool.GetTool(),
		mcp.NewStructuredToolHandler(getSectorStocksTool.HandleGetSectorStocks),
	)

	mcpServer.AddTool(
		getStockOverviewTool.GetTool(),
		mcp.NewStructuredToolHandler(getStockOverviewTool.HandleGetStockOverview),
	)

	mcpServer.AddTool(
		getStockFinancialsTool.GetTool(),
		mcp.NewStructuredToolHandler(getStockFinancialsTool.HandleGetStockFinancials),
	)

	mcpServer.AddTool(
		getUserContextTool.GetTool(),
		mcp.NewStructuredToolHandler(getUserContextTool.HandleGetUserContext),
	)

	mcpServer.AddTool(
		updateUserContextTool.GetTool(),
		mcp.NewStructuredToolHandler(updateUserContextTool.HandleUpdateUserContext),
	)

	mcpServer.AddTool(
		getEconomicIndicatorTimeSeriesTool.GetTool(),
		mcp.NewStructuredToolHandler(getEconomicIndicatorTimeSeriesTool.HandleGetEconomicIndicatorTimeSeries),
	)

	// Start the server
	httpServer := server.NewStreamableHTTPServer(mcpServer)
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
