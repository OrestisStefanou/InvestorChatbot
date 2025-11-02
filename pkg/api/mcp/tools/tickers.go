package tools

import (
	"context"
	"investbot/pkg/domain"
	"investbot/pkg/services"

	"github.com/mark3labs/mcp-go/mcp"
)

type SearchRequest struct {
	SearchString string `json:"search_string,omitempty" jsonschema_description:"Search string query"`
	Limit        int    `json:"limit,omitempty" jsonschema_description:"Maximum results" jsonschema:"minimum=1,default=100"`
}

type SearchResultSchema struct {
	Symbol      string `json:"symbol" jsonschema_description:"Stock symbol"`
	CompanyName string `json:"company_name" jsonschema_description:"Company name"`
}

type SearchResultsResponse struct {
	SearchResults []SearchResultSchema `json:"search_results" jsonschema_description:"Search results"`
}

type TickerService interface {
	GetTickers(filters services.TickerFilterOptions) ([]domain.Ticker, error)
}

type StockSearchTool struct {
	tickerService TickerService
}

func NewStockSearchTool(tickerService TickerService) (*StockSearchTool, error) {
	return &StockSearchTool{
		tickerService: tickerService,
	}, nil
}

func (h *StockSearchTool) HandleSearchStocks(ctx context.Context, req mcp.CallToolRequest, args SearchRequest) (SearchResultsResponse, error) {
	tickerFilters := services.TickerFilterOptions{
		Limit:        args.Limit,
		SearchString: args.SearchString,
	}

	tickers, err := h.tickerService.GetTickers(tickerFilters)
	if err != nil {
		return SearchResultsResponse{}, err
	}

	response := SearchResultsResponse{
		SearchResults: make([]SearchResultSchema, 0, len(tickers)),
	}

	for _, t := range tickers {
		response.SearchResults = append(
			response.SearchResults,
			SearchResultSchema{
				Symbol:      t.Symbol,
				CompanyName: t.CompanyName,
			},
		)
	}

	return response, nil
}

func (h *StockSearchTool) GetTool() mcp.Tool {
	return mcp.NewTool("stockSearch",
		mcp.WithDescription("Search for a stock using the symbol or the company name"),
		mcp.WithInputSchema[SearchRequest](),
		mcp.WithOutputSchema[SearchResultsResponse](),
	)
}
