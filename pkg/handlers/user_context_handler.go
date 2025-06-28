package handlers

import (
	"errors"
	"fmt"
	"investbot/pkg/domain"
	"net/http"

	investbotErr "investbot/pkg/errors"

	"github.com/labstack/echo/v4"
)

type UserContextService interface {
	GetUserContext(userID string) (domain.UserContext, error)
	CreateUserContext(domain.UserContext) error
	UpdateUserContext(domain.UserContext) error
}

type UserContextHandler struct {
	userContextService UserContextService
}

func NewUserContextHandler(userContextService UserContextService) (*UserContextHandler, error) {
	return &UserContextHandler{userContextService: userContextService}, nil
}

type UserPortfolioHolding struct {
	AssetClass          string  `json:"asset_class"`
	Symbol              string  `json:"symbol"`
	Name                string  `json:"name"`
	Quantity            float64 `json:"quantity"`
	PortfolioPercentage float64 `json:"portfolio_percentage"`
}

func (h UserPortfolioHolding) validate() error {
	if h.AssetClass == "" {
		return fmt.Errorf("asset_class is required")
	}

	if h.AssetClass != "stock" && h.AssetClass != "etf" && h.AssetClass != "crypto" {
		return fmt.Errorf("asset_class valid values are: stock, etf, crypto")
	}

	if h.Symbol == "" && h.Name == "" {
		return fmt.Errorf("you must define either symbol or name")
	}

	return nil
}

type UserContext struct {
	UserID        string                 `json:"user_id"`
	UserProfile   map[string]any         `json:"user_profile"`
	UserPortfolio []UserPortfolioHolding `json:"user_portfolio"`
}

func (r UserContext) validate() error {
	if r.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	for _, h := range r.UserPortfolio {
		h.validate()
	}

	return nil
}

type CreateUserContextRequest struct {
	UserContext
}

func (h *UserContextHandler) CreateUserContext(c echo.Context) error {
	request := CreateUserContextRequest{}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := request.validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	portfolioHoldings := make([]domain.UserPortfolioHolding, 0, len(request.UserPortfolio))
	for _, h := range request.UserPortfolio {
		portfolioHoldings = append(
			portfolioHoldings,
			domain.UserPortfolioHolding{
				AssetClass:          domain.AssetClass(h.AssetClass),
				Symbol:              h.Symbol,
				Name:                h.Name,
				Quantity:            h.Quantity,
				PortfolioPercentage: h.PortfolioPercentage,
			},
		)
	}

	userContext := domain.UserContext{
		UserID:        request.UserID,
		UserProfile:   request.UserProfile,
		UserPortfolio: portfolioHoldings,
	}

	err := h.userContextService.CreateUserContext(userContext)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, request)
}

type GetUserContextResponse struct {
	UserContext
}

func (h *UserContextHandler) GetUserContext(c echo.Context) error {
	userID := c.Param("user_id")

	userContext, err := h.userContextService.GetUserContext(userID)
	if err != nil {
		notFoundError := investbotErr.UserContextNotFoundError{UserID: userID}
		if errors.As(err, &notFoundError) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := GetUserContextResponse{
		UserContext: UserContext{
			UserID:        userContext.UserID,
			UserProfile:   userContext.UserProfile,
			UserPortfolio: make([]UserPortfolioHolding, 0, len(userContext.UserPortfolio)),
		},
	}

	for _, h := range userContext.UserPortfolio {
		response.UserPortfolio = append(
			response.UserPortfolio,
			UserPortfolioHolding{
				AssetClass:          string(h.AssetClass),
				Symbol:              h.Symbol,
				Name:                h.Name,
				Quantity:            h.Quantity,
				PortfolioPercentage: h.PortfolioPercentage,
			},
		)
	}

	return c.JSON(http.StatusOK, response)
}

type UpdateUserContextRequest struct {
	UserContext
}

func (h *UserContextHandler) UpdateUserContext(c echo.Context) error {
	request := UpdateUserContextRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := request.validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	portfolioHoldings := make([]domain.UserPortfolioHolding, 0, len(request.UserPortfolio))
	for _, h := range request.UserPortfolio {
		portfolioHoldings = append(
			portfolioHoldings,
			domain.UserPortfolioHolding{
				AssetClass:          domain.AssetClass(h.AssetClass),
				Symbol:              h.Symbol,
				Name:                h.Name,
				Quantity:            h.Quantity,
				PortfolioPercentage: h.PortfolioPercentage,
			},
		)
	}

	userContext := domain.UserContext{
		UserID:        request.UserID,
		UserProfile:   request.UserProfile,
		UserPortfolio: portfolioHoldings,
	}

	err := h.userContextService.UpdateUserContext(userContext)
	if err != nil {
		notFoundError := investbotErr.UserContextNotFoundError{UserID: request.UserID}
		if errors.As(err, &notFoundError) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, request)
}
