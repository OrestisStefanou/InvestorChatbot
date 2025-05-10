package domain

type RiskAppetite string

const (
	Conservative RiskAppetite = "conservative"
	Balanced     RiskAppetite = "balanced"
	Growth       RiskAppetite = "growth"
	High         RiskAppetite = "high"
)

type User struct {
	ID           string
	Email        string
	Name         string
	RiskAppetite RiskAppetite
}
