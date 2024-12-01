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
	openAiClient, _ := openAI.NewOpenAiClient(config.OpenAiKey, "https://api.openai.com/v1")

	openAiLLM, err := openAI.NewOpenAiLLM(openAI.GPT4_MINI, openAiClient, 0.2)
	if err != nil {
		log.Fatal(err)
	}

	dataService := marketDataScraper.MarketDataScraper{}
	sessionService := services.MockSessionService{}
	sectorsService := services.SectorServiceRag{
		DataService:    dataService,
		Llm:            openAiLLM,
		SessionService: sessionService,
	}
	educationService := services.EducationServiceRag{
		Llm:            openAiLLM,
		SessionService: sessionService,
	}

	topicToRagMap := make(map[services.Topic]services.Rag)
	topicToRagMap[services.EDUCATION] = educationService
	topicToRagMap[services.SECTORS] = sectorsService

	chatService, _ := services.NewChatService(topicToRagMap, nil)

	chunkChannel := make(chan string)

	go func() {
		question := "Why Do Bond Prices and Interest Rates Have an Inverse Relationship?"
		if err := chatService.GenerateResponse(services.EDUCATION, "sessionID", question, chunkChannel); err != nil {
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
