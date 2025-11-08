package tools

import (
	"context"
	"investbot/pkg/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

type SuperInvestorsService interface {
	GetSuperInvestors() ([]domain.SuperInvestor, error)
	GetSuperInvestorPortfolio(superInvestorName string) (domain.SuperInvestorPortfolio, error)
}

type SuperInvestorSchema struct {
	Name string `json:"name" jsonschema_description:"The name of the super investor(Portfolio Manager - Firm)"`
}

type GetSuperInvestorsRequest struct {
	// No input parameters required
}

type GetSuperInvestorsResponse struct {
	SuperInvestors []SuperInvestorSchema `json:"super_investors" jsonschema_description:"A list with the names of the super investors(Portfolio Managers - Firms) names"`
}

type GetSuperInvestorsTool struct {
	superInvestorsService SuperInvestorsService
}

func NewGetSuperInvestorsTool(superInvestorsService SuperInvestorsService) (*GetSuperInvestorsTool, error) {
	return &GetSuperInvestorsTool{
		superInvestorsService: superInvestorsService,
	}, nil
}

func (t *GetSuperInvestorsTool) HandleGetSuperInvestors(ctx context.Context, req mcp.CallToolRequest, args GetSuperInvestorsRequest) (GetSuperInvestorsResponse, error) {
	superInvestors, err := t.superInvestorsService.GetSuperInvestors()
	if err != nil {
		return GetSuperInvestorsResponse{}, err
	}

	response := GetSuperInvestorsResponse{
		SuperInvestors: make([]SuperInvestorSchema, 0, len(superInvestors)),
	}

	for _, si := range superInvestors {
		response.SuperInvestors = append(
			response.SuperInvestors,
			SuperInvestorSchema{Name: si.Name},
		)
	}

	return response, nil
}

func (t *GetSuperInvestorsTool) GetTool() mcp.Tool {
	return mcp.NewTool("getSuperInvestors",
		mcp.WithDescription("Get a list of all super investors (Portfolio Managers - Firms)"),
		mcp.WithInputSchema[GetSuperInvestorsRequest](),
		mcp.WithOutputSchema[GetSuperInvestorsResponse](),
	)
}
