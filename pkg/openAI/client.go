package openAI

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"investbot/pkg/errors"
	"net/http"
	"strings"
)

type chunk struct {
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

type ChatParameters struct {
	ModelName   string
	Messages    []map[string]string
	Temperature float64
}

type OpenAiClient struct {
	apiKey  string
	baseUrl string
}

func NewOpenAiClient(apiKey, baseUrl string) (*OpenAiClient, error) {
	if apiKey == "" || baseUrl == "" {
		return nil, fmt.Errorf("apiKey and baseUrl are required")
	}
	return &OpenAiClient{apiKey: apiKey, baseUrl: baseUrl}, nil
}

// Chat sends a streaming chat request to the OpenAI API, using the specified model and message history.
// It streams the response in chunks, sending each chunk's content to a provided channel for real-time processing.
//
// Parameters:
//   - parameters: A ChatParameters struct containing the model name, message history, and temperature for the chat request.
//   - chunkChannel: A channel for streaming the chat response in chunks. Each chunk of response content is sent over the
//     channel as a string. This allows for real-time processing of the response content as it arrives.
//
// Returns:
//   - error: Returns an error if any step in the process fails, including JSON marshaling, HTTP request creation,
//     network issues, or response parsing. Returns nil on success.
//
// Process:
//  1. Constructs a POST request to the OpenAI API's chat endpoint with the model and messages provided in the payload.
//  2. Initiates a streaming connection and reads chunks of the response as they arrive.
//  3. Sends each parsed chunk of content to the `chunkChannel` for real-time processing.
//  4. Closes the channel when the streaming ends or if an error occurs.
//
// Errors:
//   - Returns an error if the JSON payload cannot be marshaled, the HTTP request cannot be created,
//     the HTTP request fails, or the response contains a non-OK status code.
//   - Returns an error if JSON parsing of individual chunks fails or if an error occurs while reading the stream.
func (client OpenAiClient) Chat(parameters ChatParameters, chunkChannel chan<- string) error {
	url := fmt.Sprintf("%s/chat/completions", client.baseUrl)

	// Define the request payload
	payload := map[string]interface{}{
		"model":       parameters.ModelName,
		"messages":    parameters.Messages,
		"stream":      true,
		"temperature": parameters.Temperature,
	}

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return &errors.JSONMarshalError{
			Message: "failed to marshal JSON payload",
			Err:     err,
		}
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to create HTTP request: %v", err),
		}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.apiKey)

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return &errors.HTTPError{
			StatusCode: 0,
			Message:    fmt.Sprintf("failed to send HTTP request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return &errors.HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Create a scanner to stream the response
	scanner := bufio.NewScanner(resp.Body)
	// var chunk chunk
	for scanner.Scan() {
		chunk := chunk{}
		chunkBytes := scanner.Bytes()

		// Check if the chunk indicates the end of the stream
		if string(chunkBytes) == "data: [DONE]" {
			break
		}

		// Process only chunks that start with "data: "
		if strings.HasPrefix(string(chunkBytes), "data: ") {
			// Remove the "data: " prefix by slicing the byte slice directly
			jsonData := chunkBytes[6:] // Slice off the "data: " part (6 bytes)

			// Parse the JSON data into a Chunk struct
			if err := json.Unmarshal(jsonData, &chunk); err != nil {
				return &errors.StreamError{
					Message: "failed to parse JSON chunk",
					Err:     err,
				}
			}

			// Send the chunk content to the channel
			chunkChannel <- chunk.Choices[0].Delta.Content
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return &errors.StreamError{
			Message: "error reading the stream",
			Err:     err,
		}
	}

	close(chunkChannel)

	return nil
}
