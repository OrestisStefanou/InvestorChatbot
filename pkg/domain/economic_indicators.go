package domain

type EconomicIndicator string

const (
	RealGDP          EconomicIndicator = "RealGDP"
	TreasuryYield    EconomicIndicator = "TreasuryYield"
	InterestRate     EconomicIndicator = "InterestRate"
	Inflation        EconomicIndicator = "Inflation"
	UnemploymentRate EconomicIndicator = "UnemploymentRate"
)

type EconomicIndicatorInterval string

const (
	MonthlyEconomicIndicatorInterval   EconomicIndicatorInterval = "Monthly"
	AnnualEconomicIndicatorInterval    EconomicIndicatorInterval = "Annual"
	QuarterlyEconomicIndicatorInterval EconomicIndicatorInterval = "Quarterly"
)

type EconomicIndicatorUnit string

const (
	PercentEconomicIndicatorUnit           EconomicIndicatorUnit = "Percent"
	DollarsPerBarrelEconomicIndicatorUnit  EconomicIndicatorUnit = "DollarsPerBaller"
	BillionsOfDollarsEconomicIndicatorUnit EconomicIndicatorUnit = "BillionsOfDollars"
)

type EconomicIndicatortypeTimeSeriesEntry struct {
	Date  string
	Value string
}

type EconomicIndicatorTimeSeries struct {
	Name     EconomicIndicator
	Interval EconomicIndicatorInterval
	Unit     EconomicIndicatorUnit
	Data     []EconomicIndicatortypeTimeSeriesEntry
}

type TreasuryYieldMaturity string

const (
	ThreeMonthTreasuryYieldMaturity TreasuryYieldMaturity = "3m"
	TwoYearTreasuryYieldMaturity    TreasuryYieldMaturity = "2Y"
	FiveYearTreasuryYieldMaturity   TreasuryYieldMaturity = "5Y"
	TenYearTreasuryYieldMaturity    TreasuryYieldMaturity = "10Y"
	ThirtyYearTreasuryYieldMaturity TreasuryYieldMaturity = "30Y"
)
