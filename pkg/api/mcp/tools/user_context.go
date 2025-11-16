package tools

import (
	"context"
	"fmt"
	"investbot/pkg/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

type UserContextService interface {
	GetUserContext(userID string) (domain.UserContext, error)
	CreateUserContext(domain.UserContext) error
	UpdateUserContext(domain.UserContext) error
}

type GetUserContextRequest struct {
	UserID string `json:"user_id" jsonschema_description:"The id of the user to get the context for"`
}

type UserPortfolioHoldingSchema struct {
	AssetClass          string  `json:"asset_class" jsonschema_description:"Asset class of the holding(stock, etf, crypto etc.)"`
	Symbol              string  `json:"symbol" jsonschema_description:"Symbol of the holding"`
	Name                string  `json:"name" jsonschema_description:"Name of the holding"`
	Quantity            float64 `json:"quantity" jsonschema_description:"Quantity of the holding(zero value means not known/given)"`
	PortfolioPercentage float64 `json:"portfolio_percentage" jsonschema_description:"Portfolio percentage of the holding(zero value means not known/given)"`
}

type UserContextResponse struct {
	UserID        string                       `json:"user_id"`
	UserProfile   map[string]any               `json:"user_profile" jsonschema_description:"General information about the user"`
	UserPortfolio []UserPortfolioHoldingSchema `json:"user_portfolio"`
}

type GetUserContextTool struct {
	userContextService UserContextService
}

func NewGetUserContextTool(userContextService UserContextService) (*GetUserContextTool, error) {
	return &GetUserContextTool{
		userContextService: userContextService,
	}, nil
}

func (t *GetUserContextTool) HandleGetUserContext(ctx context.Context, req mcp.CallToolRequest, args GetUserContextRequest) (UserContextResponse, error) {
	if args.UserID == "" {
		return UserContextResponse{}, fmt.Errorf("user_id is required")
	}

	userContext, err := t.userContextService.GetUserContext(args.UserID)
	if err != nil {
		return UserContextResponse{}, err
	}

	portfolio := make([]UserPortfolioHoldingSchema, 0, len(userContext.UserPortfolio))
	for _, holding := range userContext.UserPortfolio {
		portfolio = append(portfolio, UserPortfolioHoldingSchema{
			AssetClass:          string(holding.AssetClass),
			Symbol:              holding.Symbol,
			Name:                holding.Name,
			Quantity:            holding.Quantity,
			PortfolioPercentage: holding.PortfolioPercentage,
		})
	}

	response := UserContextResponse{
		UserID:        userContext.UserID,
		UserProfile:   userContext.UserProfile,
		UserPortfolio: portfolio,
	}

	return response, nil
}

func (t *GetUserContextTool) GetTool() mcp.Tool {
	return mcp.NewTool("getUserContext",
		mcp.WithDescription("Get the user context including user profile and portfolio holdings"),
		mcp.WithInputSchema[GetUserContextRequest](),
		mcp.WithOutputSchema[UserContextResponse](),
	)
}
