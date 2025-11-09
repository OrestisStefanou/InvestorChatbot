package tools

import (
	"context"
	"investbot/pkg/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

type StockOverviewService interface {
	GetStockProfile(symbol string) (domain.StockProfile, error)
	GetFinancialRatios(symbol string) ([]domain.FinancialRatios, error)
	GetStockForecast(symbol string) (domain.StockForecast, error)
	GetHistoricalPrices(ticker string, assetClass domain.AssetClass, period domain.Period) (domain.HistoricalPrices, error)
}

type GetStockOverviewRequest struct {
	StockSymbol string `json:"stock_symbol" jsonschema_description:"Symbol of the stock to get data for"`
}

type StockProfileSchema struct {
	Name        string `json:"name" jsonschema_description:"Company name"`
	Description string `json:"description" jsonschema_description:"Company description"`
	Country     string `json:"country" jsonschema_description:"Country where company is headquartered"`
	Founded     int    `json:"founded" jsonschema_description:"Year company was founded"`
	IpoDate     string `json:"ipo_date" jsonschema_description:"Initial public offering date"`
	Industry    string `json:"industry" jsonschema_description:"Company industry"`
	Sector      string `json:"sector" jsonschema_description:"Company sector"`
	Ceo         string `json:"ceo" jsonschema_description:"Chief executive officer name"`
}

type StockEstimationSchema struct {
	Date          string  `json:"date" jsonschema_description:"Estimation date"`
	Eps           float64 `json:"eps" jsonschema_description:"Earnings per share estimate(zero value means we don't have estimation)"`
	EpsGrowth     float64 `json:"eps_growth" jsonschema_description:"EPS growth percentage(zero value means we don't have estimation)"`
	FiscalQuarter string  `json:"fiscal_quarter" jsonschema_description:"Fiscal quarter"`
	FiscalYear    string  `json:"fiscal_year" jsonschema_description:"Fiscal year"`
	Revenue       float64 `json:"revenue" jsonschema_description:"Revenue estimate(zero value means we don't have estimation)"`
	RevenueGrowth float64 `json:"revenue_growth" jsonschema_description:"Revenue growth percentage(zero value means we don't have estimation)"`
}

type StockTargetPrcSchema struct {
	Average float32 `json:"average" jsonschema_description:"Average target price(zero value means we don't have estimation)"`
	High    float32 `json:"high" jsonschema_description:"Highest target price(zero value means we don't have estimation)"`
	Low     float32 `json:"low" jsonschema_description:"Lowest target price(zero value means we don't have estimation)"`
	Median  float32 `json:"median" jsonschema_description:"Median target price(zero value means we don't have estimation)"`
}

type StockForecastSchema struct {
	Estimations []StockEstimationSchema `json:"estimations" jsonschema_description:"List of stock estimations"`
	TargetPrice StockTargetPrcSchema    `json:"target_price" jsonschema_description:"Analyst target price range"`
}

type FinancialRatiosSchema struct {
	Datekey           string  `json:"datekey" jsonschema_description:"Date key"`
	FiscalYear        string  `json:"fiscal_year" jsonschema_description:"Fiscal year"`
	FiscalQuarter     string  `json:"fiscal_quarter" jsonschema_description:"Fiscal quarter"`
	Marketcap         float64 `json:"marketcap" jsonschema_description:"Market capitalization (zero value means not known)"`
	MarketCapGrowth   float64 `json:"market_cap_growth" jsonschema_description:"Market cap growth percentage (zero value means not known)"`
	Ev                float64 `json:"ev" jsonschema_description:"Enterprise value (zero value means not known)"`
	LastCloseRatios   float64 `json:"last_close_ratios" jsonschema_description:"Last close price ratios (zero value means not known)"`
	Pe                float64 `json:"pe" jsonschema_description:"Price-to-earnings ratio (zero value means not known)"`
	Ps                float64 `json:"ps" jsonschema_description:"Price-to-sales ratio (zero value means not known)"`
	Pb                float64 `json:"pb" jsonschema_description:"Price-to-book ratio (zero value means not known)"`
	Pfcf              float64 `json:"pfcf" jsonschema_description:"Price-to-free-cash-flow ratio (zero value means not known)"`
	Pocf              float64 `json:"pocf" jsonschema_description:"Price-to-operating-cash-flow ratio (zero value means not known)"`
	EvRevenue         float64 `json:"ev_revenue" jsonschema_description:"Enterprise value to revenue ratio (zero value means not known)"`
	EvEbitda          float64 `json:"ev_ebitda" jsonschema_description:"Enterprise value to EBITDA ratio (zero value means not known)"`
	EvEbit            float64 `json:"ev_ebit" jsonschema_description:"Enterprise value to EBIT ratio (zero value means not known)"`
	EvFcf             float64 `json:"ev_fcf" jsonschema_description:"Enterprise value to free cash flow ratio (zero value means not known)"`
	DebtEquity        float64 `json:"debt_equity" jsonschema_description:"Debt-to-equity ratio (zero value means not known)"`
	DebtEbitda        float64 `json:"debt_ebitda" jsonschema_description:"Debt-to-EBITDA ratio (zero value means not known)"`
	DebtFcf           float64 `json:"debt_fcf" jsonschema_description:"Debt-to-free-cash-flow ratio (zero value means not known)"`
	AssetTurnover     float64 `json:"asset_turnover" jsonschema_description:"Asset turnover ratio (zero value means not known)"`
	InventoryTurnover float64 `json:"inventory_turnover" jsonschema_description:"Inventory turnover ratio (zero value means not known)"`
	QuickRatio        float64 `json:"quick_ratio" jsonschema_description:"Quick ratio (acid-test ratio) (zero value means not known)"`
	CurrentRatio      float64 `json:"current_ratio" jsonschema_description:"Current ratio (zero value means not known)"`
	Roe               float64 `json:"roe" jsonschema_description:"Return on equity percentage (zero value means not known)"`
	Roa               float64 `json:"roa" jsonschema_description:"Return on assets percentage (zero value means not known)"`
	Roic              float64 `json:"roic" jsonschema_description:"Return on invested capital percentage (zero value means not known)"`
	EarningsYield     float64 `json:"earnings_yield" jsonschema_description:"Earnings yield percentage (zero value means not known)"`
	FcfYield          float64 `json:"fcf_yield" jsonschema_description:"Free cash flow yield percentage (zero value means not known)"`
	DividendYield     float64 `json:"dividend_yield" jsonschema_description:"Dividend yield percentage (zero value means not known)"`
	PayoutRatio       float64 `json:"payout_ratio" jsonschema_description:"Dividend payout ratio percentage (zero value means not known)"`
	BuybackYield      float64 `json:"buyback_yield" jsonschema_description:"Share buyback yield percentage (zero value means not known)"`
	TotalReturn       float64 `json:"total_return" jsonschema_description:"Total return percentage (zero value means not known)"`
}

type PriceSchema struct {
	Date       string  `json:"date" jsonschema_description:"Price date in ISO format"`
	ClosePrice float64 `json:"close_price" jsonschema_description:"Closing price"`
}

type HistoricalPricesSchema struct {
	Period           string        `json:"period" jsonschema_description:"Time period (1d, 5d, 1m, 6m, 1y, 5y)"`
	Prices           []PriceSchema `json:"prices" jsonschema_description:"List of historical prices"`
	PercentageChange float64       `json:"percentage_change" jsonschema_description:"Percentage change over the period"`
}

type GetStockOverviewResponse struct {
	CurrentDate           string
	Symbol                string
	StockProfile          StockProfileSchema
	StockFinancialRatios  []FinancialRatiosSchema
	StockForecast         StockForecastSchema
	HistoricalPerformance []HistoricalPricesSchema
}

type GetStockOverviewTool struct {
	stockOverviewService StockOverviewService
}

func NewGetStockOverviewTool(stockOverviewService StockOverviewService) (*GetStockOverviewTool, error) {
	return &GetStockOverviewTool{
		stockOverviewService: stockOverviewService,
	}, nil
}

func (t *GetStockOverviewTool) HandleGetStockOverview(ctx context.Context, req mcp.CallToolRequest, args GetStockOverviewRequest) (GetStockOverviewResponse, error) {
	// TODO: For the forecast use only the last 2 entries of the response
	// For the financila ratios keep only the first 10
	return GetStockOverviewResponse{}, nil
}

func (t *GetStockOverviewTool) GetTool() mcp.Tool {
	return mcp.NewTool("getStockOverview",
		mcp.WithDescription("Get an overview(profile, financial rations, forecasts, performance) of the stock with the given symbol."),
		mcp.WithInputSchema[GetStockOverviewRequest](),
		mcp.WithOutputSchema[GetStockOverviewResponse](),
	)
}
