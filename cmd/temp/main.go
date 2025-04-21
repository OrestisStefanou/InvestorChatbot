package main

import (
	"fmt"
	"investbot/pkg/config"
	"investbot/pkg/marketDataScraper"
	"investbot/pkg/openAI"
	"investbot/pkg/services"
	"log"
)

func main() {
	conf, _ := config.LoadConfig()
	openAiClient, _ := openAI.NewOpenAiClient(conf.OpenAiKey, conf.OpenAiBaseUrl)
	llm, err := openAI.NewOpenAiLLM(conf.OpenAiModelName, openAiClient, float64(conf.BaseLlmTemperature))

	scraper := marketDataScraper.MarketDataScraper{}

	financialRatios, err := scraper.GetFinancialRatios("wlds")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(financialRatios)

	stockOverviewRag, err := services.NewStockOverviewRag(llm, scraper)
	if err != nil {
		log.Fatal(err)
	}

	stockProfile, err := scraper.GetStockProfile("AMD")
	if err != nil {
		log.Fatal(err)
	}

	competitors, err := stockOverviewRag.GetStockCompetitors(stockProfile, "AMD")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(competitors)
}
