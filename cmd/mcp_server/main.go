package main

import (
	"investbot/pkg/api/mcp/tools"
	"investbot/pkg/config"
	"investbot/pkg/marketDataScraper"
	"investbot/pkg/services"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

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

	// Setup cache and data services
	cache, _ := services.NewBadgerCacheService()
	dataService := marketDataScraper.NewMarketDataScraperWithCache(cache, conf)

	// Set up services
	tickerService, _ := services.NewTickerService(dataService)

	// Setup tools
	searchStocksTool, _ := tools.NewStockSearchTool(tickerService)

	// Add tools
	mcpServer.AddTool(
		searchStocksTool.GetTool(),
		mcp.NewStructuredToolHandler(searchStocksTool.HandleSearchStocks),
	)

	// Start the server
	httpServer := server.NewStreamableHTTPServer(mcpServer)
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
