package handlers

import (
	"investbot/pkg/domain"
	"investbot/pkg/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EtfService interface {
	GetEtfs(filters services.EtfFilterOptions) ([]domain.Etf, error)
}

type EtfHandler struct {
	etfService EtfService
}

type Etf struct {
	Symbol     string  `json:"symbol"`
	Name       string  `json:"name"`
	AssetClass string  `json:"asset_class"`
	Aum        float32 `json:"aum"`
}

type GetEtfsResponse struct {
	Etfs []Etf `json:"etfs"`
}

func NewEtfHandler(etfService EtfService) (*EtfHandler, error) {
	return &EtfHandler{etfService: etfService}, nil
}

func (h *EtfHandler) GetEtfs(c echo.Context) error {
	searchString := c.QueryParam("search_string")

	tickerFilters := services.EtfFilterOptions{
		SearchString: searchString,
	}

	etfs, err := h.etfService.GetEtfs(tickerFilters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := GetEtfsResponse{
		Etfs: make([]Etf, 0, len(etfs)),
	}

	for _, e := range etfs {
		response.Etfs = append(
			response.Etfs,
			Etf{Symbol: e.Symbol, Name: e.Name, AssetClass: e.AssetClass, Aum: e.Aum},
		)
	}

	return c.JSON(http.StatusOK, response)
}
