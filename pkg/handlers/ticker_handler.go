package handlers

import (
	"investbot/pkg/domain"
	"investbot/pkg/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TickerService interface {
	GetTickers(filters services.TickerFilterOptions) ([]domain.Ticker, error)
}

type TickerHandler struct {
	tickerService TickerService
}

type Ticker struct {
	Symbol      string `json:"symbol"`
	CompanyName string `json:"company_name"`
}

type GetTickersResponse struct {
	Tickers []Ticker `json:"tickers"`
}

func NewTickerHandler(tickerService TickerService) (*TickerHandler, error) {
	return &TickerHandler{
		tickerService: tickerService,
	}, nil
}

func (h *TickerHandler) GetTickers(c echo.Context) error {
	// Get filter parameters
	var err error
	var limit, page int

	limitQueryParam := c.QueryParam("limit")
	if limitQueryParam == "" {
		limit = 0
	} else {
		limit, err = strconv.Atoi(limitQueryParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "limit query param must be a valid integer"})
		}
	}

	pageQueryParam := c.QueryParam("page")
	if pageQueryParam == "" {
		page = 0
	} else {
		page, err = strconv.Atoi(pageQueryParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "page query param must be a valid integer"})
		}
	}

	searchString := c.QueryParam("search_string")

	tickerFilters := services.TickerFilterOptions{
		Limit:        limit,
		Page:         page,
		SearchString: searchString,
	}

	tickers, err := h.tickerService.GetTickers(tickerFilters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := GetTickersResponse{
		Tickers: make([]Ticker, 0, len(tickers)),
	}

	for _, t := range tickers {
		response.Tickers = append(response.Tickers, Ticker{Symbol: t.Symbol, CompanyName: t.CompanyName})
	}

	return c.JSON(http.StatusOK, response)
}
