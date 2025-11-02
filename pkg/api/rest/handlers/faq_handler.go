package handlers

import (
	"investbot/pkg/errors"
	"investbot/pkg/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FaqService interface {
	GetFaqForTopic(topic services.FaqTopic) ([]string, error)
}

type FaqHandler struct {
	faqService FaqService
}

func NewFaqHandler(faqService FaqService) (*FaqHandler, error) {
	return &FaqHandler{
		faqService: faqService,
	}, nil
}

type GetFaqResponse struct {
	Faq []string `json:"faq"`
}

func (h *FaqHandler) GetFaq(c echo.Context) error {
	// Get topic from query parameter
	topic := c.QueryParam("faq_topic")
	faq, err := h.faqService.GetFaqForTopic(services.FaqTopic(topic))

	if err != nil {
		switch e := err.(type) {
		case *errors.FaqTopicNotFoundError:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": e.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

	}

	response := GetFaqResponse{Faq: faq}
	return c.JSON(http.StatusOK, response)
}
