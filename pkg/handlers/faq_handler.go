package handlers

import (
	"investbot/pkg/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FaqHandler struct {
	faqService services.FaqService
}

func NewFaqHandler(faqService services.FaqService) (*FaqHandler, error) {
	return &FaqHandler{
		faqService: faqService,
	}, nil
}

type GetFaqResponse struct {
	Faq []string `json:"faq"`
}

func (h *FaqHandler) GetFaq(c echo.Context) error {
	// Get topic from query parameter
	topic := c.QueryParam("topic")
	faq, err := h.faqService.GetFaqForTopic(services.Topic(topic))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	response := GetFaqResponse{Faq: faq}
	return c.JSON(http.StatusOK, response)
}
