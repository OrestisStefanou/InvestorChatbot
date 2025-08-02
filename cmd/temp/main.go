package main

import (
	"fmt"
	"investbot/pkg/config"
	"investbot/pkg/domain"
	"investbot/pkg/llama"
	"investbot/pkg/marketDataScraper"
	"investbot/pkg/openAI"
	"investbot/pkg/services"
	"log"
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
	conf, _ := config.LoadConfig()
	cache, _ := services.NewBadgerCacheService()
	dataService := marketDataScraper.NewMarketDataScraperWithCache(cache, conf)

	historicalPrices, err := dataService.GetHistoricalPrices("VOO", domain.ETF, domain.Period6M)
	if err != nil {
		log.Fatal(err)
	}

	// historicalPrices, err = dataService.GetHistoricalPrices("meta", domain.Stock, domain.Period5Y)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Printf("\n\n%+v\n\n", historicalPrices)
}
