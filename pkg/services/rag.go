package services

import "strings"

type RagResponsesStore interface {
	StoreRagResponse(
		ragTopic Topic,
		prompt string,
		conversation []Message,
		response string,
	) error
}

type BaseRag struct {
	topic         Topic
	llm           Llm
	responseStore RagResponsesStore
}

func (r *BaseRag) GenerateLllmResponse(
	prompt string,
	conversation []Message,
	responseChannel chan<- string,
) error {
	// Step 1: Add prompt as first user message
	conversation = append([]Message{{Content: prompt, Role: User}}, conversation...)

	// Step 2: Buffered chunks to prevent writer blocking
	chunks := make(chan string, 16)
	defer close(chunks) // guarantees goroutine exit

	var fullResponse strings.Builder

	// Step 3: Forward chunks and build full response
	go func() {
		defer close(responseChannel)
		for chunk := range chunks {
			responseChannel <- chunk
			fullResponse.WriteString(chunk)
		}
	}()

	// Step 4: Call LLM to generate chunks
	if err := r.llm.GenerateResponse(conversation, chunks); err != nil {
		return err
	}

	// Step 5: Store full response
	return r.responseStore.StoreRagResponse(
		r.topic,
		prompt,
		conversation,
		fullResponse.String(),
	)
}
