package domain

import "time"

type AssetClass string

const (
	Stock  AssetClass = "stock"
	ETF    AssetClass = "etf"
	Crypto AssetClass = "crypto"
)

type RiskLevel string

const (
	LowRiskLevel    RiskLevel = "low"
	MediumRiskLevel RiskLevel = "medium"
	HighRiskLevel   RiskLevel = "high"
)

type PortfolioHolding struct {
	AssetClass AssetClass
	AssetID    string // This is usually the symbol of the asset
	Quantity   float64
}

type Portfolio struct {
	ID        string
	Name      string
	RiskLevel RiskLevel
	Holdings  []PortfolioHolding
	CreatedAt time.Time
	UpdatedAt time.Time
}
