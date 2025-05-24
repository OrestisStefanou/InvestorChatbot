package handlers

import (
	"errors"
	"fmt"
	"investbot/pkg/domain"
	"net/http"
	"time"

	investbotErr "investbot/pkg/errors"

	"github.com/labstack/echo/v4"
)

type PortfolioService interface {
	GetPortfolioById(portfolioID string) (domain.Portfolio, error)
	CreatePortfolio(portfolio domain.Portfolio) error
	UpdatePortfolio(portfolio domain.Portfolio) error
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
	PortfolioID string             `json:"portfolio_id"`
	Name        string             `json:"name"`
	RiskLevel   string             `json:"risk_level"`
	Holdings    []PortfolioHolding `json:"holdings"`
}

func (r CreatePortfolioRequest) validate() error {
	if r.PortfolioID == "" {
		return fmt.Errorf("portfolio_id is required")
	}

	if r.RiskLevel != "" {
		if r.RiskLevel != "low" && r.RiskLevel != "medium" && r.RiskLevel != "high" {
			return fmt.Errorf("risk_level valid values are: low, medium, high")
		}
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

func (h *PortfolioHandler) CreatePortfolio(c echo.Context) error {
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
		ID:        request.PortfolioID,
		Name:      request.Name,
		RiskLevel: domain.RiskLevel(request.RiskLevel),
		Holdings:  portfolioHoldings,
	}

	err := h.portfolioService.CreatePortfolio(portfolio)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"Message": "Portfolio created"})
}

type GetPortfolioResponse struct {
	Name      string             `json:"name"`
	RiskLevel string             `json:"risk_level"`
	Holdings  []PortfolioHolding `json:"holdings"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
}

func (h *PortfolioHandler) GetPortfolioById(c echo.Context) error {
	portfolioID := c.Param("portfolio_id")

	portfolio, err := h.portfolioService.GetPortfolioById(portfolioID)
	if err != nil {
		notFoundError := investbotErr.PortfolioNotFoundError{PortfolioID: portfolioID}
		if errors.As(err, &notFoundError) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := GetPortfolioResponse{
		Name:      portfolio.Name,
		RiskLevel: string(portfolio.RiskLevel),
		Holdings:  make([]PortfolioHolding, 0, len(portfolio.Holdings)),
		CreatedAt: portfolio.CreatedAt.Format(time.RFC3339),
		UpdatedAt: portfolio.UpdatedAt.Format(time.RFC3339),
	}
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

type UpdatePortfolioRequest struct {
	Name      string             `json:"name"`
	RiskLevel string             `json:"risk_level"`
	Holdings  []PortfolioHolding `json:"holdings"`
}

func (h *PortfolioHandler) UpdatePortfolio(c echo.Context) error {
	portfolioID := c.Param("portfolio_id")

	request := UpdatePortfolioRequest{}
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
		ID:        portfolioID,
		Name:      request.Name,
		RiskLevel: domain.RiskLevel(request.RiskLevel),
		Holdings:  portfolioHoldings,
	}

	err := h.portfolioService.UpdatePortfolio(portfolio)
	if err != nil {
		notFoundError := investbotErr.PortfolioNotFoundError{PortfolioID: portfolioID}
		if errors.As(err, &notFoundError) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"Message": "Portfolio updated"})
}
