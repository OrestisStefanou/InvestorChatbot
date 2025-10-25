package handlers

import (
	"fmt"
	"investbot/pkg/errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FollowUpService interface {
	GenerateFollowUpQuestions(sessionId string, followUpQuestionsNum int) ([]string, error)
}

type FollowUpQuestionsHandler struct {
	followUpQuestionsService FollowUpService
	followUpQuestionsNum     int
}

type FollowUpQuestionsResponse struct {
	FollowUpQuestions []string `json:"follow_up_questions"`
}

type FollowUpQuestionsRequest struct {
	SessionID         string `json:"session_id"`
	NumberOfQuestions int    `json:"number_of_questions"`
}

func (r *FollowUpQuestionsRequest) validate() error {
	if r.SessionID == "" {
		return fmt.Errorf("session_id field is required")
	}

	return nil
}

func NewFollowUpQuestionsHandler(followUpQuestionsService FollowUpService, followUpQuestionsNum int) (*FollowUpQuestionsHandler, error) {
	return &FollowUpQuestionsHandler{
		followUpQuestionsService: followUpQuestionsService,
		followUpQuestionsNum:     followUpQuestionsNum,
	}, nil
}

func (h *FollowUpQuestionsHandler) GenerateFollowUpQuestions(c echo.Context) error {
	var err error
	request := new(FollowUpQuestionsRequest)
	if err = c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = request.validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if request.NumberOfQuestions == 0 {
		request.NumberOfQuestions = h.followUpQuestionsNum
	}

	followUpQuestions, err := h.followUpQuestionsService.GenerateFollowUpQuestions(request.SessionID, request.NumberOfQuestions)

	if err != nil {
		switch e := err.(type) {
		case *errors.SessionNotFoundError:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": e.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	response := FollowUpQuestionsResponse{FollowUpQuestions: followUpQuestions}
	return c.JSON(http.StatusOK, response)
}
