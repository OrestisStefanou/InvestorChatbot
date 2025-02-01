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
	chunkChannel := make(chan string)

	go func() {
		messages := []services.Message{
			{Role: services.User, Content: "Which are the top performing sectors?"},
		}
		if err := sectorRag.GenerateRagResponse(messages, chunkChannel); err != nil {
			// Handle the error (e.g., log it)
			log.Printf("Error during request: %v", err)
		}
	}()
	// Consume the chunks from the channel
	var finalResponse string
	for content := range chunkChannel {
		// Process the chunk as it arrives
		fmt.Printf("Received chunk: %s\n", content)
		finalResponse += content
	}

	// Optional: After the channel is closed, perform any final tasks
	log.Println("Streaming has finished.")
	fmt.Println(finalResponse)
}
