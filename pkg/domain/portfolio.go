package domain

type AssetClass string

const (
	Stock  AssetClass = "stock"
	ETF    AssetClass = "etf"
	Crypto AssetClass = "crypto"
)

type PortfolioHolding struct {
	AssetClass AssetClass
	AssetID    string // This is usually the symbol of the asset
	Quantity   float64
}

type Portfolio struct {
	UserEmail string
	Holdings  []PortfolioHolding
}
