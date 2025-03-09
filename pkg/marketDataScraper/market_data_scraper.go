package marketDataScraper

import "investbot/pkg/domain"

type MarketDataScraper struct {
}

// GetSectorStocks returns a list of stocks in a sector
// sector parameter should be the domain.Sector.UrlName value
func (mds MarketDataScraper) GetSectorStocks(sector string) ([]domain.SectorStock, error) {
	return scrapeSectorStocks(sector)
}

// GetSectors returns a list of sectors
func (mds MarketDataScraper) GetSectors() ([]domain.Sector, error) {
	return scrapeSectors()
}

// GetIndustryStocks returns a list of stocks in an industry
// industry parameter should be the domain.Industry.UrlName value
func (mds MarketDataScraper) GetIndustryStocks(industry string) ([]domain.IndustryStock, error) {
	return scrapeIndustryStocks(industry)
}

// GetIndustries returns a list of industries
func (mds MarketDataScraper) GetIndustries() ([]domain.Industry, error) {
	return scrapeIndustries()
}

// GetStockForecsat returns the forecast for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetStockForecast(symbol string) (domain.StockForecast, error) {
	return scrapeStockForecast(symbol)
}

// GetBalanceSheets returns a list of balance sheets for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetBalanceSheets(symbol string) ([]domain.BalanceSheet, error) {
	return scrapeBalanceSheets(symbol)
}

// GetIncomeStatements returns a list of income statements for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetIncomeStatements(symbol string) ([]domain.IncomeStatement, error) {
	return scrapeIncomeStatements(symbol)
}

// GetCashFlows returns a list of cash flows for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetCashFlows(symbol string) ([]domain.CashFlow, error) {
	return scrapeCashFlows(symbol)
}

// GetFinancialRatios returns a list of financial ratios for a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetFinancialRatios(symbol string) ([]domain.FinancialRatios, error) {
	return scrapeFinancialRatios(symbol)
}

// GetEtfs returns a list of ETFs
func (mds MarketDataScraper) GetEtfs() ([]domain.Etf, error) {
	return scrapeEtfs()
}

// GetEtfOverview returns an overview of an ETF
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetEtfOverview(symbol string) (domain.EtfOverview, error) {
	return scrapeEtfOverview(symbol)
}

// GetStockProfile returns the profile of a stock
// symbol parameter should be in lowercase
func (mds MarketDataScraper) GetStockProfile(symbol string) (domain.StockProfile, error) {
	return scrapeStockProfile(symbol)
}

// GetMarketNews returns the most recent news of the stock markets
func (mds MarketDataScraper) GetMarketNews() ([]domain.NewsArticle, error) {
	return scrapeMarketNews()
}
