package handlers

import (
	"investbot/pkg/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TopicHandler struct{}

type GetTopicsResponse struct {
	Topics []string `json:"topics"`
}

func NewTopicHandler() (*TopicHandler, error) {
	return &TopicHandler{}, nil
}

func (h *TopicHandler) GetTopics(c echo.Context) error {
	topics := []string{
		string(services.EDUCATION),
		string(services.SECTORS),
		string(services.STOCK_OVERVIEW),
		string(services.STOCK_FINANCIALS),
		string(services.ETFS),
		string(services.NEWS),
	}

	response := GetTopicsResponse{Topics: topics}

	return c.JSON(http.StatusOK, response)
}
