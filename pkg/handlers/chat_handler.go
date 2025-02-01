package handlers

import (
	"encoding/json"
	"fmt"
	"investbot/pkg/services"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChatService interface {
	GenerateResponse(topic services.Topic, sessionId string, question string, responseChannel chan<- string) error
}

type ChatHandler struct {
	chatService ChatService
}

func NewChatHandler(chatService ChatService) (*ChatHandler, error) {
	return &ChatHandler{chatService: chatService}, nil
}

type ChatRequest struct {
	Question  string `json:"question" validate:"required"`
	Topic     string `json:"topic" validate:"required"`
	SessionID string `json:"session_id" validate:"required"`
	// TODO: Add extra optional data(metadata) in the request
	// For example stock symbol/name etc
}

func (h *ChatHandler) ServeRequest(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)

	var err error
	chatRequest := new(ChatRequest)
	if err = c.Bind(chatRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	enc := json.NewEncoder(c.Response())
	responseChunkChannel := make(chan string)

	go func() {
		if err = h.chatService.GenerateResponse(
			services.Topic(chatRequest.Topic), chatRequest.SessionID, chatRequest.Question, responseChunkChannel,
		); err != nil {
			// TODO: Think which is the right way to handle the error
			log.Println("ERROR DURING GENERATE RESPONSE", err)
			return
		}
	}()

	var finalResponse string
	for chunk := range responseChunkChannel {
		if err := enc.Encode(chunk); err != nil {
			return err
		}
		finalResponse += chunk
		c.Response().Flush()
	}
	fmt.Printf("FINAL RESPONSE\n %s", finalResponse)
	return nil

}
