package llama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"investbot/pkg/errors"
	"net/http"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chunk struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Message   Message   `json:"message"`
	Done      bool      `json:"done"`
}

type Options struct {
	Temperature float32 `json:"temperature"`
}

type ChatParameters struct {
	ModelName string    `json:"model"`
	Messages  []Message `json:"messages"`
	Options   Options   `json:"options"`
	Stream    bool      `json:"stream"`
}

type OllamaClient struct {
	baseUrl string
}

func NewOllamaClient(baseUrl string) (*OllamaClient, error) {
	return &OllamaClient{baseUrl: baseUrl}, nil
}

func (client *OllamaClient) Chat(parameters ChatParameters, chunkChannel chan<- string) error {
	defer close(chunkChannel)
	url := fmt.Sprintf("%s/api/chat", client.baseUrl)

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(parameters)
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
	var chunk chunk
	for scanner.Scan() {
		chunkBytes := scanner.Bytes()

		if err := json.Unmarshal(chunkBytes, &chunk); err != nil {
			return &errors.StreamError{
				Message: "failed to parse JSON chunk",
				Err:     err,
			}
		}

		chunkChannel <- chunk.Message.Content

		// Check if the chunk indicates the end of the stream
		if chunk.Done {
			break
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return &errors.StreamError{
			Message: "error reading the stream",
			Err:     err,
		}
	}

	return nil
}
