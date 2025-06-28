package domain

type UserPortfolioHolding struct {
	AssetClass          AssetClass
	Symbol              string
	Name                string
	Quantity            float64
	PortfolioPercentage float64
}

type UserContext struct {
	UserID        string
	UserProfile   map[string]any
	UserPortfolio []UserPortfolioHolding
}
