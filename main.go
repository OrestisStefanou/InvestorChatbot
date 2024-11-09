package main

import (
	"fmt"
	"investbot/openAiClient"
	"log"
)

func main() {
	// Create a channel to receive chunk data
	chunkChannel := make(chan string)

	// Start the MakeChatRequest function in a goroutine to stream data
	go func() {
		if err := openAiClient.MakeChatRequest(chunkChannel); err != nil {
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
