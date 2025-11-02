package handlers

import (
	investbotErr "investbot/pkg/errors"
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

type Actor string

const (
	Assistant Actor = "assistant"
	User      Actor = "user"
)

type Message struct {
	Actor   Actor  `json:"actor"`
	Message string `json:"message"`
}

type GetSessionResponse struct {
	Conversation []Message `json:"conversation"`
}

func (h *SessionHandler) GetSession(c echo.Context) error {
	sessionID := c.Param("session_id")

	session, err := h.sessionService.GetConversationBySessionId(sessionID)
	if err != nil {
		switch e := err.(type) {
		case *investbotErr.SessionNotFoundError:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": e.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": e.Error()})
		}
	}

	response := GetSessionResponse{}
	response.Conversation = make([]Message, 0, len(session))
	for _, m := range session {
		var actor Actor
		switch m.Role {
		case services.Assistant:
			actor = Assistant
		case services.User:
			actor = User
		default:
			continue
		}

		msg := Message{Actor: actor, Message: m.Content}
		response.Conversation = append(response.Conversation, msg)
	}

	return c.JSON(http.StatusOK, response)
}
