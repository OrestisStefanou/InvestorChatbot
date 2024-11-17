package main

import (
	"fmt"
	"investbot/pkg/marketDataScraper"
	"investbot/pkg/openAI"
	"investbot/pkg/services"
	"log"
)

func main() {
	openAiClient, _ := openAI.NewOpenAiClient("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", "https://api.openai.com/v1")

	openAiLLM := openAI.OpenAiLLM{
		ModelName:   "gpt-4o-mini",
		Client:      openAiClient,
		Temperature: 0.2,
	}
	// Start the MakeChatRequest function in a goroutine to stream data
	dataService := marketDataScraper.MarketDataScraper{}
	sessionService := services.MockSessionService{}
	sectorService := services.SectorServiceRag{
		DataService:    dataService,
		Llm:            openAiLLM,
		SessionService: sessionService,
	}
	chunkChannel := make(chan string)

	go func() {
		// conversation := []map[string]string{
		// 	{"role": "user", "content": "Hey there!"},
		// 	{"role": "system", "content": "Hello! How can I help you today?"},
		// 	{"role": "user", "content": "What is a synonym for big?"},
		// }
		sessionId := "test_session_id"
		question := "Which stocks are in the Technology sector?"
		if err := sectorService.GenerateRagResponse(sessionId, question, chunkChannel); err != nil {
			// Handle the error (e.g., log it)
			log.Printf("Error during request: %v", err)
			close(chunkChannel) // Ensure the channel is closed if there’s an error
		}
	}()

	// Consume the chunks from the channel
	for content := range chunkChannel {
		// Process the chunk as it arrives
		fmt.Printf("Received chunk: %s\n", content)
	}

	// Optional: After the channel is closed, perform any final tasks
	log.Println("Streaming has finished.")

}
