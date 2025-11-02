package handlers

import (
	"errors"
	"investbot/pkg/domain"
	investbotErr "investbot/pkg/errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SuperInvestorService interface {
	GetSuperInvestors() ([]domain.SuperInvestor, error)
	GetSuperInvestorPortfolio(superInvestorName string) (domain.SuperInvestorPortfolio, error)
}

type SuperInvestorHandler struct {
	superInvestorService SuperInvestorService
}

type SuperInvestorPortfolioHolding struct {
	Stock          string `json:"stock"`
	PortfolioPct   string `json:"portfolio_pct"`
	RecentActivity string `json:"recent_activity"`
	Shares         string `json:"shares"`
	Value          string `json:"value"`
}

type SuperInvestorPortfolioSectorAnalysis struct {
	Sector       string `json:"sector"`
	PortfolioPct string `json:"portfolio_pct"`
}

type SuperInvestor struct {
	Name string `json:"name"`
}

type GetSuperInvestorPortfolioResponse struct {
	Holdings       []SuperInvestorPortfolioHolding        `json:"holdings"`
	SectorAnalysis []SuperInvestorPortfolioSectorAnalysis `json:"sector_analysis"`
}

type GetSuperInvestorsResponse struct {
	SuperInvestors []SuperInvestor `json:"super_investors"`
}

func NewSuperInvestorHandler(superInvestorService SuperInvestorService) (*SuperInvestorHandler, error) {
	return &SuperInvestorHandler{superInvestorService: superInvestorService}, nil
}

func (h *SuperInvestorHandler) GetSuperInvestors(c echo.Context) error {
	superInvestors, err := h.superInvestorService.GetSuperInvestors()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := GetSuperInvestorsResponse{
		SuperInvestors: make([]SuperInvestor, 0, len(superInvestors)),
	}

	for _, s := range superInvestors {
		response.SuperInvestors = append(response.SuperInvestors, SuperInvestor{Name: s.Name})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *SuperInvestorHandler) GetSuperInvestorPortfolio(c echo.Context) error {
	superInvestorName := c.Param("super_investor")
	portfolio, err := h.superInvestorService.GetSuperInvestorPortfolio(superInvestorName)
	if err != nil {
		notFoundError := investbotErr.SuperInvestorPortfolioNotFoundError{}
		if errors.As(err, &notFoundError) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := GetSuperInvestorPortfolioResponse{
		Holdings:       make([]SuperInvestorPortfolioHolding, 0, len(portfolio.Holdings)),
		SectorAnalysis: make([]SuperInvestorPortfolioSectorAnalysis, 0, len(portfolio.SectorAnalysis)),
	}

	for _, h := range portfolio.Holdings {
		response.Holdings = append(
			response.Holdings,
			SuperInvestorPortfolioHolding{
				Stock:          h.Stock,
				PortfolioPct:   h.PortfolioPct,
				RecentActivity: h.RecentActivity,
				Shares:         h.Shares,
				Value:          h.Value,
			},
		)
	}

	for _, s := range portfolio.SectorAnalysis {
		response.SectorAnalysis = append(
			response.SectorAnalysis,
			SuperInvestorPortfolioSectorAnalysis{
				Sector:       s.Sector,
				PortfolioPct: s.PortfolioPct,
			},
		)
	}

	return c.JSON(http.StatusOK, response)
}
