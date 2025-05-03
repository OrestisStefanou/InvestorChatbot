package handlers

import (
	"errors"
	"investbot/pkg/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	GetUser(id string) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) (*UserHandler, error) {
	return &UserHandler{userService: userService}, nil
}

type CreateUserRequest struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	RiskAppetite string `json:"risk_appetite"`
}

func (r *CreateUserRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}

	if r.Name == "" {
		return errors.New("name is required")
	}

	if r.RiskAppetite != "" {
		if r.RiskAppetite != "conservative" && r.RiskAppetite != "balanced" && r.RiskAppetite != "growth" && r.RiskAppetite != "high" {
			return errors.New("invalid risk appetite, valid values are: conservative, balanced, growth, high")
		}
	}

	return nil
}

type CreateUserResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	RiskAppetite string `json:"risk_appetite"`
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	request := CreateUserRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, CreateUserResponse{ID: ""})
	}

	if err := request.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user := domain.User{
		Email:        request.Email,
		Name:         request.Name,
		RiskAppetite: domain.RiskAppetite(request.RiskAppetite),
	}

	userCreated, err := h.userService.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, CreateUserResponse{
		ID:           userCreated.ID,
		Email:        userCreated.Email,
		Name:         userCreated.Name,
		RiskAppetite: string(userCreated.RiskAppetite),
	})
}
