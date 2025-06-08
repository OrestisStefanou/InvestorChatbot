package handlers

import (
	"encoding/json"
	"fmt"
	investbotErr "investbot/pkg/errors"
	"investbot/pkg/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChatService interface {
	GenerateResponse(topic services.Topic, tags services.Tags, sessionId string, question string, responseChannel chan<- string) error
	ExtractTopicAndTags(question string, sessionId string) (services.Topic, services.Tags, error)
}

type ChatHandler struct {
	chatService ChatService
}

func NewChatHandler(chatService ChatService) (*ChatHandler, error) {
	return &ChatHandler{chatService: chatService}, nil
}

type TopicTags struct {
	SectorName      string   `json:"sector_name"`
	IndustryName    string   `json:"industry_name"`
	StockSymbols    []string `json:"stock_symbols"`
	BalanceSheet    bool     `json:"balance_sheet"`
	IncomeStatement bool     `json:"income_statement"`
	CashFlow        bool     `json:"cash_flow"`
	EtfSymbol       string   `json:"etf_symbol"`
}

type ChatRequest struct {
	Question  string    `json:"question"`
	Topic     string    `json:"topic"`
	SessionID string    `json:"session_id"`
	Tags      TopicTags `json:"topic_tags"`
}

func (r *ChatRequest) validate() error {
	if r.Question == "" {
		return fmt.Errorf("question field is required")
	}

	if r.Topic == "" {
		return fmt.Errorf("topic field is required")
	}

	if r.SessionID == "" {
		return fmt.Errorf("session_id field is required")
	}

	return nil
}

func (h *ChatHandler) ChatCompletion(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var err error
	chatRequest := new(ChatRequest)
	if err = c.Bind(chatRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = chatRequest.validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tags := services.Tags{
		SectorName:      chatRequest.Tags.SectorName,
		IndustryName:    chatRequest.Tags.IndustryName,
		StockSymbols:    chatRequest.Tags.StockSymbols,
		BalanceSheet:    chatRequest.Tags.BalanceSheet,
		IncomeStatement: chatRequest.Tags.IncomeStatement,
		CashFlow:        chatRequest.Tags.CashFlow,
		EtfSymbol:       chatRequest.Tags.EtfSymbol,
	}

	enc := json.NewEncoder(c.Response())
	responseChunkChannel := make(chan string)
	errorChannel := make(chan error, 1)

	go func() {
		if err := h.chatService.GenerateResponse(
			services.Topic(chatRequest.Topic), tags, chatRequest.SessionID, chatRequest.Question, responseChunkChannel,
		); err != nil {
			errorChannel <- err
		}
		close(errorChannel)
	}()

	for {
		select {
		case chunk, isOpen := <-responseChunkChannel:
			if !isOpen {
				// Channel closed, exit loop
				return nil
			}
			if err = enc.Encode(chunk); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}

			c.Response().WriteHeader(http.StatusOK)
			c.Response().Flush()

		case err := <-errorChannel:
			if err != nil {
				switch e := err.(type) {
				case *investbotErr.SessionNotFoundError:
					return c.JSON(http.StatusBadRequest, map[string]string{"error": e.Error()})
				case *investbotErr.InvalidTopicError:
					return c.JSON(http.StatusBadRequest, map[string]string{"error": e.Error()})
				default:
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": e.Error()})
				}

			}
		}
	}
}

type ExtractTopicAndTagsRequest struct {
	Question  string `json:"question"`
	SessionID string `json:"session_id"`
}

func (r ExtractTopicAndTagsRequest) validate() error {
	if r.Question == "" {
		return fmt.Errorf("question field is required")
	}

	if r.SessionID == "" {
		return fmt.Errorf("session_id field is required")
	}

	return nil
}

type ExtractTopicAndTagsResponse struct {
	Topic string    `json:"topic"`
	Tags  TopicTags `json:"topic_tags"`
}

func (h *ChatHandler) ExtractTopicAndTags(c echo.Context) error {
	extractTopicAndTagsRequest := new(ExtractTopicAndTagsRequest)
	if err := c.Bind(extractTopicAndTagsRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := extractTopicAndTagsRequest.validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	topic, tags, err := h.chatService.ExtractTopicAndTags(extractTopicAndTagsRequest.Question, extractTopicAndTagsRequest.SessionID)
	if err != nil {
		switch e := err.(type) {
		case *investbotErr.SessionNotFoundError:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": e.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	response := ExtractTopicAndTagsResponse{
		Topic: string(topic),
		Tags: TopicTags{
			SectorName:      tags.SectorName,
			IndustryName:    tags.IndustryName,
			StockSymbols:    tags.StockSymbols,
			BalanceSheet:    tags.BalanceSheet,
			IncomeStatement: tags.IncomeStatement,
			CashFlow:        tags.CashFlow,
			EtfSymbol:       tags.EtfSymbol,
		},
	}

	return c.JSON(http.StatusOK, response)
}
