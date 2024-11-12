package domain

type Etf struct {
	Symbol     string
	Name       string
	AssetClass string
	Aum        float32
}

type EtfHolding struct {
	Symbol string
	Name   string
	Weight string
}

type EtfOverview struct {
	Symbol           string
	Description      string
	AssetClass       string
	Category         string
	Aum              string
	Nav              string
	ExpenseRatio     string
	PeRatio          string
	Dps              string
	DividendYield    string
	PayoutRatio      string
	OneYearReturn    string
	Beta             string
	NumberOfHoldings int32
	Website          string
	TopHoldings      []EtfHolding
}
