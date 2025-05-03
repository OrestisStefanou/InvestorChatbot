package domain

type AssetClass string

const (
	Stock  AssetClass = "stock"
	ETF    AssetClass = "etf"
	Crypto AssetClass = "crypto"
)

type PortfolioHolding struct {
	AssetClass AssetClass
	AssetID    string
	Quantity   float64
}

type Portfolio struct {
	ID       string
	UserID   string
	Holdings []PortfolioHolding
}
