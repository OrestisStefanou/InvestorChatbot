package handlers

import (
	"investbot/pkg/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SectorService interface {
	GetSectorStocks(sector string) ([]domain.SectorStock, error)
	GetSectors() ([]domain.Sector, error)
}

type SectorHandler struct {
	sectorService SectorService
}

type Sector struct {
	Name             string  `json:"name"`
	UrlName          string  `json:"url_name"`
	NumberOfStocks   int     `json:"number_of_stocks"`
	MarketCap        float32 `json:"market_cap"`
	DividendYieldPct float32 `json:"dividend_yield_pct"`
	PeRatio          float32 `json:"pe_ratio"`
	ProfitMarginPct  float32 `json:"profit_margin_pct"`
	OneYearChangePct float32 `json:"one_year_change_pct"`
}

type SectorStock struct {
	Symbol      string  `json:"symbol"`
	CompanyName string  `json:"company_name"`
	MarketCap   float32 `json:"market_cap"`
}

type GetSectorsResponse struct {
	Sectors []Sector `json:"sectors"`
}

type GetSectorStocksResponse struct {
	SectorStocks []SectorStock `json:"sector_stocks"`
}

func NewSectorHandler(sectorService SectorService) (*SectorHandler, error) {
	return &SectorHandler{sectorService: sectorService}, nil
}

func (h *SectorHandler) GetSectors(c echo.Context) error {
	sectors, err := h.sectorService.GetSectors()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := GetSectorsResponse{
		Sectors: make([]Sector, 0, len(sectors)),
	}

	for _, s := range sectors {
		response.Sectors = append(
			response.Sectors,
			Sector{
				Name:             s.Name,
				UrlName:          s.UrlName,
				NumberOfStocks:   s.NumberOfStocks,
				MarketCap:        s.MarketCap,
				DividendYieldPct: s.DividendYieldPct,
				PeRatio:          s.PeRatio,
				ProfitMarginPct:  s.ProfitMarginPct,
				OneYearChangePct: s.OneYearChangePct,
			},
		)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *SectorHandler) GetSectorStocks(c echo.Context) error {
	sector := c.Param("sector")
	sectorStocks, err := h.sectorService.GetSectorStocks(sector)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := GetSectorStocksResponse{
		SectorStocks: make([]SectorStock, 0, len(sectorStocks)),
	}

	for _, s := range sectorStocks {
		response.SectorStocks = append(
			response.SectorStocks,
			SectorStock{
				Symbol:      s.Symbol,
				CompanyName: s.CompanyName,
				MarketCap:   s.MarketCap,
			},
		)
	}

	return c.JSON(http.StatusOK, response)
}
