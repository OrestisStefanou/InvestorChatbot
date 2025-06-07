package main

import (
	"fmt"
	"investbot/pkg/config"
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
	llm, err := getLlm(conf)
	if err != nil {
		log.Fatal(err)
	}
	cache, _ := services.NewBadgerCacheService()
	dataService := marketDataScraper.NewMarketDataScraperWithCache(cache, conf)
	// topicExtractor, _ := services.NewTopicExtractor(llm)
	// conversation := []services.Message{
	// 	services.Message{
	// 		Content: "What are the benefits of investing?",
	// 		Role:    services.User,
	// 	},
	// 	services.Message{
	// 		Content: "It helps you grow your money over time to beat inflation",
	// 		Role:    services.Assistant,
	// 	},
	// 	services.Message{
	// 		Content: "What is inflation?",
	// 		Role:    services.User,
	// 	},
	// }
	// topic, err := topicExtractor.ExtractTopic(conversation)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Topic is: %s\n", topic)
	tagExtractor, _ := services.NewTagExtractor(llm, dataService)
	conversation := []services.Message{
		services.Message{
			Content: "Why should I invest in etfs?",
			Role:    services.User,
		},
		// services.Message{
		// 	Content: "It helps you grow your money over time to beat inflation",
		// 	Role:    services.Assistant,
		// },
		// services.Message{
		// 	Content: "What is inflation?",
		// 	Role:    services.User,
		// },
	}
	tags, err := tagExtractor.ExtractTags(services.ETFS, conversation)
	if err != nil {
		fmt.Printf("Err: %s", err)
	}
	fmt.Printf("Tags: %+v", tags)
}
