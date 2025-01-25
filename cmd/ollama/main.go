package main

import (
	"fmt"
	"investbot/pkg/llama"
	"log"
)

func main() {
	llamaClient, _ := llama.NewOllamaClient("http://localhost:11434")

	chunkChannel := make(chan string)

	go func() {
		parameters := llama.ChatParameters{
			ModelName: "llama3.2",
			Messages: []llama.Message{
				{Role: "user", Content: "Hey there!"},
			},
			Options: llama.Options{
				Temperature: 0.1,
			},
			Stream: true,
		}
		if err := llamaClient.Chat(parameters, chunkChannel); err != nil {
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
