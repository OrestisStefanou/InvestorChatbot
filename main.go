package main

import (
	"fmt"
	"investbot/openAI"
	"log"
)

func main() {
	// Create a channel to receive chunk data
	chunkChannel := make(chan string)

	openAiClient, _ := openAI.NewOpenAiClient("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", "https://api.openai.com/v1")
	// Start the MakeChatRequest function in a goroutine to stream data
	go func() {
		parameters := openAI.ChatParameters{
			ModelName:   "gpt-4o-mini",
			Temperature: 0.5,
			Messages: []map[string]string{
				{"role": "system", "content": "You are a helpful assistant."},
				{"role": "user", "content": "Hey there!"},
				{"role": "system", "content": "Hello! How can I help you today?"},
				{"role": "user", "content": "What is a synonym for big?"},
			},
		}
		if err := openAiClient.Chat(parameters, chunkChannel); err != nil {
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
