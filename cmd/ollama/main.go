package main

import (
	"fmt"
	"investbot/pkg/config"
	"investbot/pkg/gemini"
	"investbot/pkg/services"
	"log"
)

func main() {
	conf, _ := config.LoadConfig()
	llmConfig := gemini.GeminiLlmConfig{
		ModelName:   gemini.GEMINI_2_0_FLASH,
		Temperature: 0.1,
		ApiKey:      conf.GeminiKey,
	}
	llm, _ := gemini.NewGeminiLLM(llmConfig)

	messages := []services.Message{
		{Role: services.System, Content: "You are an investing expert named Warren Buffet jr."},
		{Role: services.User, Content: "Hey there, what's your name?"},
	}

	responseChannel := make(chan string)
	errorChannel := make(chan error, 1)
	go func() {
		err := llm.GenerateResponse(messages, responseChannel)
		if err != nil {
			errorChannel <- err
		}
		close(errorChannel)
	}()

	shouldExit := false
	var responseMessage string
	for !shouldExit {
		select {
		case chunk, isOpen := <-responseChannel:
			if !isOpen {
				fmt.Printf("\n\nFINAL RESPONSE\n %s", responseMessage)
				shouldExit = true
				continue
			}
			responseMessage += chunk
		case err := <-errorChannel:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
