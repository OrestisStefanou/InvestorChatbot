package handlers

import (
	"fmt"
	"investbot/pkg/domain"
	"net/http"

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

		if h.AssetClass != "stock" || h.AssetClass != "etf" || h.AssetClass != "crypto" {
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
