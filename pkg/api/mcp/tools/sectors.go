package tools

import (
	"context"
	"investbot/pkg/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

type SectorSchema struct {
	Name             string  `json:"name" jsonschema_description:"Sector name"`
	UrlName          string  `json:"url_name" jsonschema_description:"Used for internal purposes"`
	NumberOfStocks   int     `json:"number_of_stocks" jsonschema_description:"Number of stocks in the sector"`
	MarketCap        float32 `json:"market_cap" jsonschema_description:"Market cap of the secotr"`
	DividendYieldPct float32 `json:"dividend_yield_pct" jsonschema_description:"Dividend yield percentage of the sector"`
	PeRatio          float32 `json:"pe_ratio" jsonschema_description:"PE ratio of the sector"`
	ProfitMarginPct  float32 `json:"profit_margin_pct" jsonschema_description:"Profit margin percentage of the sector"`
	OneYearChangePct float32 `json:"one_year_change_pct" jsonschema_description:"One year price change percentage of the sector"`
}

type SectorStockSchema struct {
	Symbol      string  `json:"symbol" jsonschema_description:"Stock symbol"`
	CompanyName string  `json:"company_name" jsonschema_description:"Company name of the stock"`
	MarketCap   float32 `json:"market_cap" jsonschema_description:"Market cap of the stock"`
}

// TODO:
//  Add a tool to return the stocks of a sector

type SectorsService interface {
	GetSectorStocks(sector string) ([]domain.SectorStock, error)
	GetSectors() ([]domain.Sector, error)
}

type GetSectorsRequest struct {
	// No input parameters required
}

type GetSectorsResponse struct {
	Sectors []SectorSchema `json:"sectors" jsonschema_description:"A list with the sectors"`
}

type GetSectorsTool struct {
	sectorsService SectorsService
}

func NewGetSectorsTool(sectorsService SectorsService) (*GetSectorsTool, error) {
	return &GetSectorsTool{
		sectorsService: sectorsService,
	}, nil
}

func (t *GetSectorsTool) HandleGetSectors(ctx context.Context, req mcp.CallToolRequest, args GetSectorsRequest) (GetSectorsResponse, error) {
	sectors, err := t.sectorsService.GetSectors()
	if err != nil {
		return GetSectorsResponse{}, err
	}

	response := GetSectorsResponse{Sectors: make([]SectorSchema, 0, len(sectors))}

	for _, sector := range sectors {
		response.Sectors = append(response.Sectors, SectorSchema{
			Name:             sector.Name,
			UrlName:          sector.UrlName,
			NumberOfStocks:   sector.NumberOfStocks,
			MarketCap:        sector.MarketCap,
			DividendYieldPct: sector.DividendYieldPct,
			PeRatio:          sector.PeRatio,
			ProfitMarginPct:  sector.ProfitMarginPct,
			OneYearChangePct: sector.OneYearChangePct,
		})
	}

	return response, nil
}

func (t *GetSectorsTool) GetTool() mcp.Tool {
	return mcp.NewTool("getSectors",
		mcp.WithDescription("Get all stock sectors"),
		mcp.WithInputSchema[GetSectorsRequest](),
		mcp.WithOutputSchema[GetSectorsResponse](),
	)
}
