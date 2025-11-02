package tools

import (
	"context"
	"investbot/pkg/domain"
	"investbot/pkg/services"

	"github.com/mark3labs/mcp-go/mcp"
)

type SearchEtfRequest struct {
	SearchString string `json:"search_string,omitempty" jsonschema_description:"Search string query"`
	Limit        int    `json:"limit,omitempty" jsonschema_description:"Maximum results" jsonschema:"minimum=1,default=100"`
}

type EtfSearchResultSchema struct {
	Symbol     string  `json:"symbol" jsonschema_description:"ETF symbol"`
	EtfName    string  `json:"etf_name" jsonschema_description:"ETF name"`
	AssetClass string  `json:"asset_class" jsonschema_description:"ETF asset class"`
	Aum        float32 `json:"aum" jsonschema_description:"ETF assets under management"`
}

type EtfSearchResultsResponse struct {
	SearchResults []EtfSearchResultSchema `json:"search_results" jsonschema_description:"Search results"`
}

type EtfService interface {
	GetEtfs(filters services.EtfFilterOptions) ([]domain.Etf, error)
}

type SearchEtfTool struct {
	etfService EtfService
}

func NewSearchEtfTool(etfService EtfService) (*SearchEtfTool, error) {
	return &SearchEtfTool{
		etfService: etfService,
	}, nil
}

func (h *SearchEtfTool) HandleSearchEtfs(ctx context.Context, req mcp.CallToolRequest, args SearchEtfRequest) (EtfSearchResultsResponse, error) {
	if args.Limit == 0 {
		args.Limit = 100
	}
	tickerFilters := services.EtfFilterOptions{
		SearchString: args.SearchString,
	}

	etfs, err := h.etfService.GetEtfs(tickerFilters)
	if err != nil {
		return EtfSearchResultsResponse{}, err
	}

	response := EtfSearchResultsResponse{
		SearchResults: make([]EtfSearchResultSchema, 0, len(etfs)),
	}

	for i, e := range etfs {
		if i > args.Limit {
			break
		}
		response.SearchResults = append(
			response.SearchResults,
			EtfSearchResultSchema{Symbol: e.Symbol, EtfName: e.Name, AssetClass: e.AssetClass, Aum: e.Aum},
		)
	}

	return response, nil
}

func (h *SearchEtfTool) GetTool() mcp.Tool {
	return mcp.NewTool("etfSearch",
		mcp.WithDescription("Search for an ETF using the symbol or the ETF name"),
		mcp.WithInputSchema[SearchEtfRequest](),
		mcp.WithOutputSchema[EtfSearchResultsResponse](),
	)
}
