package services

import "investbot/pkg/domain"

type MarketDataService interface {
	GetSectorStocks(sector string) ([]domain.SectorStock, error)
	GetSectors() ([]domain.Sector, error)
	GetIndustryStocks(industry string) ([]domain.IndustryStock, error)
	GetIndustries() ([]domain.Industry, error)
	GetStockForecast(symbol string) (domain.StockForecast, error)
	GetBalanceSheets(symbol string) ([]domain.BalanceSheet, error)
	GetIncomeStatements(symbol string) ([]domain.IncomeStatement, error)
	GetCashFlows(symbol string) ([]domain.CashFlow, error)
	GetFinancialRatios(symbol string) ([]domain.FinancialRatios, error)
	GetEtfs() ([]domain.Etf, error)
	GetEtfOverview(symbol string) (domain.EtfOverview, error)
	GetStockProfile(symbol string) (domain.StockProfile, error)
}
