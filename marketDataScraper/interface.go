package marketDataScraper

import "investbot/domain"

// GetSectorStocks returns a list of stocks in a sector
// sector parameter should be the domain.Sector.UrlName value
func GetSectorStocks(sector string) ([]domain.SectorStock, error) {
	return scrapeSectorStocks(sector)
}

// GetSectors returns a list of sectors
func GetSectors() ([]domain.Sector, error) {
	return scrapeSectors()
}

// GetIndustryStocks returns a list of stocks in an industry
// industry parameter should be the domain.Industry.UrlName value
func GetIndustryStocks(industry string) ([]domain.IndustryStock, error) {
	return scrapeIndustryStocks(industry)
}

// GetIndustries returns a list of industries
func GetIndustries() ([]domain.Industry, error) {
	return scrapeIndustries()
}

// GetStockForecsat returns the forecast for a stock
// symbol parameter should be in lowercase
func GetStockForecast(symbol string) (domain.StockForecast, error) {
	return scrapeStockForecast(symbol)
}

// GetBalanceSheets returns a list of balance sheets for a stock
// symbol parameter should be in lowercase
func GetBalanceSheets(symbol string) ([]domain.BalanceSheet, error) {
	return scrapeBalanceSheets(symbol)
}

// GetIncomeStatements returns a list of income statements for a stock
// symbol parameter should be in lowercase
func GetIncomeStatements(symbol string) ([]domain.IncomeStatement, error) {
	return scrapeIncomeStatements(symbol)
}

// GetCashFlows returns a list of cash flows for a stock
// symbol parameter should be in lowercase
func GetCashFlows(symbol string) ([]domain.CashFlow, error) {
	return scrapeCashFlows(symbol)
}

// GetFinancialRatios returns a list of financial ratios for a stock
// symbol parameter should be in lowercase
func GetFinancialRatios(symbol string) ([]domain.FinancialRatios, error) {
	return scrapeFinancialRatios(symbol)
}

// GetEtfs returns a list of ETFs
func GetEtfs() ([]domain.Etf, error) {
	return scrapeEtfs()
}

// GetEtfOverview returns an overview of an ETF
// symbol parameter should be in lowercase
func GetEtfOverview(symbol string) (domain.EtfOverview, error) {
	return scrapeEtfOverview(symbol)
}

// GetStockProfile returns the profile of a stock
// symbol parameter should be in lowercase
func GetStockProfile(symbol string) (domain.StockProfile, error) {
	return scrapeStockProfile(symbol)
}
