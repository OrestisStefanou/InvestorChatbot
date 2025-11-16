package alphavantage

type EconomicIndicator string

const (
	RealGDP          EconomicIndicator = "REAL_GDP"
	TreasuryYield    EconomicIndicator = "TREASURY_YIELD"
	FederalFundsRate EconomicIndicator = "FEDERAL_FUNDS_RATE"
	Inflation        EconomicIndicator = "INFLATION"
	UnemploymentRate EconomicIndicator = "UNEMPLOYMENT"
)

type TimeSeriesEntry struct {
	Date  string `json:"date"`
	Value string `json:"value"`
}

type EconomicIndicatorTimeSeriesResponse struct {
	Name     string            `json:"name"`
	Interval string            `json:"interval"`
	Unit     string            `json:"unit"`
	Data     []TimeSeriesEntry `json:"data"`
}
