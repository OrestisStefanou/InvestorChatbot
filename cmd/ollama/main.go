package main

import (
	"fmt"
	"investbot/pkg/llama"
	"investbot/pkg/marketDataScraper"
	"investbot/pkg/services"
	"log"
)

func main() {
	// config, _ := config.LoadConfig()
	// openAiClient, _ := openAI.NewOpenAiClient(config.OpenAiKey, "https://api.openai.com/v1")

	// openAiLLM, err := openAI.NewOpenAiLLM(openAI.GPT4_MINI, openAiClient, 0.2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	llamaClient, _ := llama.NewOllamaClient("http://localhost:11434")
	llamaLLM, _ := llama.NewLlamaLLM("llama3.2", llamaClient, 0.2)
	dataService := marketDataScraper.MarketDataScraper{}
	sectorRag, _ := services.NewSectorRag(llamaLLM, dataService)
	chunkChannel := make(chan string)

	go func() {
		messages := []services.Message{
			{Role: services.User, Content: "Who is the best football player in the world?"},
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
