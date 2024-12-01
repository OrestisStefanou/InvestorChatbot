package openAI

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test for successful streaming
func TestChat_Success(t *testing.T) {
	// Create a mock server to simulate the OpenAI API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the method and headers
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("expected Authorization header to be 'Bearer test-api-key'")
		}

		// Stream mock response chunks
		flusher, _ := w.(http.Flusher)
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"Hello\"}}]}\n\n"))
		flusher.Flush()
		time.Sleep(10 * time.Millisecond)
		w.Write([]byte("data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"World\"}}]}\n\n"))
		flusher.Flush()
		time.Sleep(10 * time.Millisecond)
		w.Write([]byte("data: [DONE]\n\n"))
		flusher.Flush()
	}))
	defer server.Close()

	client := OpenAiClient{apiKey: "test-api-key", baseUrl: server.URL}
	chunkChannel := make(chan string)
	//defer close(chunkChannel)

	go func() {
		parameters := chatParameters{
			ModelName:   "test-model",
			Temperature: 0.5,
			Messages:    []map[string]string{{"role": "user", "content": "Hello"}},
		}
		err := client.Chat(parameters, chunkChannel)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}()

	// Collect streamed chunks and check them
	expectedChunks := []string{"Hello", "World"}
	for i, expected := range expectedChunks {
		select {
		case chunk := <-chunkChannel:
			if chunk != expected {
				t.Errorf("expected chunk %d to be '%s', got '%s'", i, expected, chunk)
			}
		case <-time.After(time.Second):
			t.Errorf("timeout waiting for chunk %d", i)
		}
	}
}

// TestChat_ErrorCases tests various error scenarios in the Chat method
func TestChat_ErrorCases(t *testing.T) {
	tests := []struct {
		name           string
		client         OpenAiClient
		messages       []map[string]string
		mockServerFunc func() *httptest.Server
		expectedError  error
	}{
		{
			name:     "HTTP Client Error",
			client:   OpenAiClient{baseUrl: "http://invalid-url"}, // invalid URL to induce client error
			messages: []map[string]string{{"role": "user", "content": "hello"}},
			mockServerFunc: func() *httptest.Server {
				return httptest.NewServer(nil)
			},
			expectedError: &HTTPError{StatusCode: http.StatusNotFound, Message: "failed to send HTTP request"},
		},
		{
			name:     "HTTP Response Error (Non-200 Status)",
			client:   OpenAiClient{baseUrl: "http://mock.url"},
			messages: []map[string]string{{"role": "user", "content": "hello"}},
			mockServerFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError) // Respond with 500 error
				}))
			},
			expectedError: &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error"},
		},
		{
			name:     "Stream Parse Error",
			client:   OpenAiClient{baseUrl: "http://mock.url"},
			messages: []map[string]string{{"role": "user", "content": "hello"}},
			mockServerFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("data: {invalid json}")) // send invalid JSON to induce parse error
				}))
			},
			expectedError: &StreamError{Message: "stream JSON parse error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock server if needed
			server := tt.mockServerFunc()
			defer server.Close()

			// Set client base URL to mock server's URL
			if server != nil {
				tt.client.baseUrl = server.URL
			}

			// Define a dummy channel to receive stream content
			chunkChannel := make(chan string)

			// Execute the Chat method
			parameters := chatParameters{
				ModelName:   "test-model",
				Temperature: 0.5,
				Messages:    tt.messages,
			}
			err := tt.client.Chat(parameters, chunkChannel)

			// Compare errors
			if err == nil {
				t.Errorf("expected error but got nil")
				return
			}

			// Type assert and check error message for each error type
			switch expectedErr := tt.expectedError.(type) {
			case *JSONMarshalError:
				if _, ok := err.(*JSONMarshalError); !ok {
					t.Errorf("expected JSONMarshalError but got %v", err)
				}
			case *HTTPError:
				if httpErr, ok := err.(*HTTPError); !ok || httpErr.StatusCode != expectedErr.StatusCode {
					t.Errorf("expected HTTPError with status %d but got %v", expectedErr.StatusCode, err)
				}
			case *StreamError:
				if _, ok := err.(*StreamError); !ok {
					t.Errorf("expected StreamError but got %v", err)
				}
			default:
				t.Errorf("unexpected error type: %v", err)
			}
		})
	}
}
