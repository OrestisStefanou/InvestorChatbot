package gemini

import (
	"context"
	"investbot/pkg/services"

	"google.golang.org/genai"
)

type ModelName string

const (
	GEMINI_2_5_FLASH ModelName = "gemini-2.5-flash"
	GEMINI_2_0_FLASH ModelName = "gemini-2.0-flash"
)

type GeminiLlmConfig struct {
	ModelName   ModelName
	Temperature float32
	ApiKey      string
}

type GeminiLLM struct {
	config GeminiLlmConfig
}

func NewGeminiLLM(llmConfig GeminiLlmConfig) (*GeminiLLM, error) {
	return &GeminiLLM{config: llmConfig}, nil
}

func (llm GeminiLLM) GenerateResponse(conversation []services.Message, responseChannel chan<- string) error {
	generateContentConfig := &genai.GenerateContentConfig{
		Temperature: &llm.config.Temperature,
	}

	conversationLen := len(conversation)
	messages := make([]*genai.Content, 0, conversationLen)
	var role genai.Role
	for _, m := range conversation[:conversationLen-1] {
		switch m.Role {
		case services.Assistant:
			role = genai.RoleModel
		case services.User:
			role = genai.RoleUser
		case services.System:
			generateContentConfig.SystemInstruction = genai.NewContentFromText(m.Content, genai.RoleUser)
			continue
		}
		messages = append(messages, genai.NewContentFromText(m.Content, role))
	}

	clientConfig := genai.ClientConfig{
		APIKey: llm.config.ApiKey,
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &clientConfig)
	if err != nil {
		return err
	}

	chat, _ := client.Chats.Create(ctx, string(llm.config.ModelName), generateContentConfig, messages)
	lastMessage := conversation[conversationLen-1].Content
	stream := chat.SendMessageStream(ctx, genai.Part{Text: lastMessage})

	for chunk, _ := range stream {
		part := chunk.Candidates[0].Content.Parts[0]
		responseChannel <- part.Text
	}

	close(responseChannel)

	return nil
}
