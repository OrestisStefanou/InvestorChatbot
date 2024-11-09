package openAiClient

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Chunk struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int64  `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Index int `json:"index"`
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason interface{} `json:"finish_reason"`
	} `json:"choices"`
}

// MakeChatRequest streams data using a channel and returns any errors encountered
func MakeChatRequest(chunkChannel chan<- string) error {
	url := "https://api.openai.com/v1/chat/completions"
	apiKey := "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"

	// Define the request payload
	payload := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant."},
			{"role": "user", "content": "Hey there!"},
			{"role": "system", "content": "Hello! How can I help you today?"},
			{"role": "user", "content": "What is a synonym for big?"},
		},
		"stream": true,
	}

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from server: %s", resp.Status)
	}

	// Create a scanner to stream the response
	scanner := bufio.NewScanner(resp.Body)
	var chunk Chunk
	for scanner.Scan() {
		chunkBytes := scanner.Bytes()

		// Check if the chunk indicates the end of the stream
		if string(chunkBytes) == "data: [DONE]" {
			close(chunkChannel) // Close the channel when the stream ends
			break
		}

		// Process only chunks that start with "data: "
		if strings.HasPrefix(string(chunkBytes), "data: ") {
			// Remove the "data: " prefix by slicing the byte slice directly
			jsonData := chunkBytes[6:] // Slice off the "data: " part (6 bytes)

			// Parse the JSON data into a Chunk struct
			if err := json.Unmarshal(jsonData, &chunk); err != nil {
				return fmt.Errorf("failed to parse JSON chunk: %w", err)
			}

			// Send the chunk content to the channel
			chunkChannel <- chunk.Choices[0].Delta.Content
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading the stream: %w", err)
	}

	return nil
}
