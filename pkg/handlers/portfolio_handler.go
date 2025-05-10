package handlers

import (
	"errors"
	"fmt"
	"investbot/pkg/domain"
	"net/http"

	investbotErr "investbot/pkg/errors"

	"github.com/labstack/echo/v4"
)

type PortfolioService interface {
	GetUserPortfolio(userEmail string) (domain.Portfolio, error)
	CreateUserPortfolio(portfolio domain.Portfolio) error
	UpdateUserPortfolio(portfolio domain.Portfolio) error
}

type PortfolioHandler struct {
	portfolioService PortfolioService
}

func NewPortfolioHandler(portfolioService PortfolioService) (*PortfolioHandler, error) {
	return &PortfolioHandler{portfolioService: portfolioService}, nil
}

type PortfolioHolding struct {
	AssetClass string  `json:"asset_class"`
	Symbol     string  `json:"symbol"`
	Quantity   float64 `json:"quantity"`
}

type CreatePortfolioRequest struct {
	UserEmail string             `json:"user_email"`
	Holdings  []PortfolioHolding `json:"holdings"`
}

func (r CreatePortfolioRequest) validate() error {
	if r.UserEmail == "" {
		return fmt.Errorf("user_email is required")
	}

	for _, h := range r.Holdings {
		if h.AssetClass == "" {
			return fmt.Errorf("asset_class is required")
		}

		if h.AssetClass != "stock" && h.AssetClass != "etf" && h.AssetClass != "crypto" {
			return fmt.Errorf("asset_class valid values are: stock, etf, crypto")
		}

		if h.Symbol == "" {
			return fmt.Errorf("symbol is required")
		}
	}

	return nil
}

func (h *PortfolioHandler) CreateUserPortfolio(c echo.Context) error {
	request := CreatePortfolioRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := request.validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	portfolioHoldings := make([]domain.PortfolioHolding, 0, len(request.Holdings))
	for _, h := range request.Holdings {
		portfolioHoldings = append(
			portfolioHoldings,
			domain.PortfolioHolding{
				AssetClass: domain.AssetClass(h.AssetClass),
				AssetID:    h.Symbol,
				Quantity:   h.Quantity,
			},
		)
	}

	portfolio := domain.Portfolio{
		UserEmail: request.UserEmail,
		Holdings:  portfolioHoldings,
	}

	err := h.portfolioService.CreateUserPortfolio(portfolio)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"Message": "Portfolio created"})
}

type GetUserPortfolioResponse struct {
	Holdings []PortfolioHolding `json:"holdings"`
}

func (h *PortfolioHandler) GetUserPortfolio(c echo.Context) error {
	userEmail := c.QueryParam("user_email") // TODO: We should extract the email from a jwt

	portfolio, err := h.portfolioService.GetUserPortfolio(userEmail)
	if err != nil {
		notFoundError := investbotErr.PortfolioNotFoundError{UserEmail: userEmail}
		if errors.As(err, &notFoundError) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := GetUserPortfolioResponse{}
	response.Holdings = make([]PortfolioHolding, 0, len(portfolio.Holdings))
	for _, h := range portfolio.Holdings {
		response.Holdings = append(
			response.Holdings,
			PortfolioHolding{
				AssetClass: string(h.AssetClass),
				Symbol:     h.AssetID,
				Quantity:   h.Quantity,
			},
		)
	}

	return c.JSON(http.StatusOK, response)
}

type UpdateUserPortfolioRequest struct {
	Holdings []PortfolioHolding `json:"holdings"`
}

func (h *PortfolioHandler) UpdateUserPortfolio(c echo.Context) error {
	userEmail := c.QueryParam("user_email") // TODO: We should extract the email from a jwt

	request := UpdateUserPortfolioRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	portfolioHoldings := make([]domain.PortfolioHolding, 0, len(request.Holdings))
	for _, h := range request.Holdings {
		portfolioHoldings = append(
			portfolioHoldings,
			domain.PortfolioHolding{
				AssetClass: domain.AssetClass(h.AssetClass),
				AssetID:    h.Symbol,
				Quantity:   h.Quantity,
			},
		)
	}

	portfolio := domain.Portfolio{
		UserEmail: userEmail,
		Holdings:  portfolioHoldings,
	}

	err := h.portfolioService.UpdateUserPortfolio(portfolio)
	if err != nil {
		notFoundError := investbotErr.PortfolioNotFoundError{UserEmail: userEmail}
		if errors.As(err, &notFoundError) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"Message": "Portfolio updated"})
}
