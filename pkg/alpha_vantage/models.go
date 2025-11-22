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

type Commodity string

const (
	WTI         Commodity = "WTI"
	NATURAL_GAS Commodity = "NATURAL_GAS"
	COPPER      Commodity = "COPPER"
	ALUMINIUM   Commodity = "ALUMINUM"
	WHEAT       Commodity = "WHEAT"
	CORN        Commodity = "CORN"
	SUGAR       Commodity = "SUGAR"
	COFFEE      Commodity = "COFFEE"
)

type CommodityTimeSeriesResponse struct {
	Name     string            `json:"name"`
	Interval string            `json:"interval"`
	Unit     string            `json:"unit"`
	Data     []TimeSeriesEntry `json:"data"`
}
