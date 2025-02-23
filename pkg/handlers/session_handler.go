package handlers

import (
	"investbot/pkg/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	sessionService services.SessionService
}

type CreateSessionResponse struct {
	SessionId string `json:"session_id"`
}

func NewSessionHandler(sessionService services.SessionService) (*SessionHandler, error) {
	return &SessionHandler{
		sessionService: sessionService,
	}, nil
}

func (h *SessionHandler) CreateNewSession(c echo.Context) error {
	sessionId, err := h.sessionService.CreateNewSession()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	response := CreateSessionResponse{SessionId: sessionId}
	return c.JSON(http.StatusCreated, response)
}
