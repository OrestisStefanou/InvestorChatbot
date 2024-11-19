package main

import (
	"fmt"
	"investbot/pkg/config"
	"investbot/pkg/openAI"
	"investbot/pkg/services"
	"log"
)

func main() {
	config, _ := config.LoadConfig()
	openAiClient, _ := openAI.NewOpenAiClient(config.OpenAiKey, "https://api.openai.com/v1")

	openAiLLM := openAI.OpenAiLLM{
		ModelName:   "gpt-4o-mini",
		Client:      openAiClient,
		Temperature: 0.2,
	}

	//dataService := marketDataScraper.MarketDataScraper{}
	sessionService := services.MockSessionService{}
	educationService := services.EducationServiceRag{
		//DataService:    dataService,
		Llm:            openAiLLM,
		SessionService: sessionService,
	}
	chunkChannel := make(chan string)

	go func() {
		question := "Why Do Bond Prices and Interest Rates Have an Inverse Relationship?"
		if err := educationService.GenerateRagResponse("session_id", question, chunkChannel); err != nil {
			// Handle the error (e.g., log it)
			log.Printf("Error during request: %v", err)
			close(chunkChannel) // Ensure the channel is closed if there’s an error
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
