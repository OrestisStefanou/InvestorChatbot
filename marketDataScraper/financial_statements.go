package marketDataScraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type BalanceSheet struct {
	Datekey                        string
	FiscalYear                     string
	FiscalQuarter                  string
	Cashneq                        float64
	Investmentsc                   float64
	Totalcash                      float64
	CashGrowth                     float64
	AccountsReceivable             float64
	OtherReceivables               float64
	Receivables                    float64
	Inventory                      float64
	RestrictedCash                 float64
	Othercurrent                   float64
	Assetsc                        float64
	NetPPE                         float64
	Investmentsnc                  float64
	Goodwill                       float64
	OtherIntangibles               float64
	Othernoncurrent                float64
	Assets                         float64
	AccountsPayable                float64
	AccruedExpenses                float64
	Debtc                          float64
	CurrentPortDebt                float64
	CurrentCapLeases               float64
	CurrentIncomeTaxesPayable      float64
	CurrentUnearnedRevenue         float64
	OtherCurrentLiabilities        float64
	CurrentLiabilities             float64
	Debtnc                         float64
	CapitalLeases                  float64
	LongTermUnearnedRevenue        float64
	LongTermDeferredTaxLiabilities float64
	Otherliabilitiesnoncurrent     float64
	Liabilities                    float64
	CommonStock                    float64
	Retearn                        float64
	OtherEquity                    float64
	Equity                         float64
	Liabilitiesequity              float64
	SharesOutFilingDate            float64
	SharesOutTotalCommon           float64
	Bvps                           float64
	TangibleBookValue              float64
	TangibleBookValuePerShare      float64
	Debt                           float64
	Netcash                        float64
	NetCashGrowth                  float64
	Netcashpershare                float64
	Workingcapital                 float64
	Land                           float64
	Machinery                      float64
	LeaseholdImprovements          float64
	TradingAssetSecurities         float64
}

type CashFlow struct {
	Datekey                  string
	FiscalYear               string
	FiscalQuarter            string
	NetIncomeCF              float64
	TotalDepAmorCF           float64
	Sbcomp                   float64
	ChangeAR                 float64
	ChangeInventory          float64
	ChangeAP                 float64
	ChangeUnearnedRev        float64
	ChangeOtherNetOperAssets float64
	OtherOperating           float64
	Ncfo                     float64
	OcfGrowth                float64
	Capex                    float64
	CashAcquisition          float64
	SalePurchaseIntangibles  float64
	InvestInSecurities       float64
	OtherInvesting           float64
	Ncfi                     float64
	DebtIssuedShortTerm      float64
	DebtIssuedLongTerm       float64
	DebtIssuedTotal          float64
	DebtRepaidShortTerm      float64
	DebtRepaidLongTerm       float64
	DebtRepaidTotal          float64
	NetDebtIssued            float64
	CommonIssued             float64
	CommonRepurchased        float64
	CommonDividendCF         float64
	OtherFinancing           float64
	Ncff                     float64
	Ncf                      float64
	Fcf                      float64
	FcfGrowth                float64
	FcfMargin                float64
	Fcfps                    float64
	LeveredFCF               float64
	UnleveredFCF             float64
	CashInterestPaid         float64
	CashTaxesPaid            float64
	ChangeNetWorkingCapital  float64
}

type IncomeStatement struct {
	Datekey           string
	FiscalYear        string
	FiscalQuarter     string
	Revenue           float64
	RevenueGrowth     float64
	Cor               float64
	Gp                float64
	Sgna              float64
	Rnd               float64
	Opex              float64
	Opinc             float64
	InterestExpense   float64
	InterestIncome    float64
	CurrencyGains     float64
	OtherNonOperating float64
	EbtExcl           float64
	GainInvestments   float64
	Pretax            float64
	Taxexp            float64
	Netinc            float64
	Netinccmn         float64
	NetIncomeGrowth   float64
	SharesBasic       float64
	SharesDiluted     float64
	SharesYoY         float64
	EpsBasic          float64
	EpsDil            float64
	EpsGrowth         float64
	Fcf               float64
	Fcfps             float64
	Dps               float64
	DividendGrowth    float64
	GrossMargin       float64
	OperatingMargin   float64
	ProfitMargin      float64
	FcfMargin         float64
	Taxrate           float64
	Ebitda            float64
	DepAmorEbitda     float64
	EbitdaMargin      float64
	Ebit              float64
	EbitMargin        float64
	RevenueAsReported float64
	PayoutRatio       float64
}

type FinancialRatios struct {
	Datekey           string
	FiscalYear        string
	FiscalQuarter     string
	Marketcap         float64
	MarketCapGrowth   float64
	Ev                float64
	LastCloseRatios   float64
	Pe                float64
	Ps                float64
	Pb                float64
	Pfcf              float64
	Pocf              float64
	EvRevenue         float64
	EvEbitda          float64
	EvEbit            float64
	EvFcf             float64
	DebtEquity        float64
	DebtEbitda        float64
	DebtFcf           float64
	AssetTurnover     float64
	InventoryTurnover float64
	QuickRatio        float64
	CurrentRatio      float64
	Roe               float64
	Roa               float64
	Roic              float64
	EarningsYield     float64
	FcfYield          float64
	DividendYield     float64
	PayoutRatio       float64
	BuybackYield      float64
	TotalReturn       float64
}

func scrapeFinancialStatementData(url string) ([]map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return []map[string]interface{}{}, err
	}

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		return []map[string]interface{}{}, errors.New("unexpected structure in 'nodes'")
	}

	// Access the second element in "nodes" which contains the data we are interested in
	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}, errors.New("unexpected structure in 'nodes[2]'")
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, errors.New("unexpected structure in 'data'")
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}, errors.New("unexpected structure in 'data[0]'")
	}

	financialDataIndex, ok := dataMap["financialData"].(float64)
	if !ok {
		return []map[string]interface{}{}, errors.New("unexpected structure for 'financialData'")
	}

	// Retrieve data map at financial data index
	statementDataMap, ok := data[int(financialDataIndex)].(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}, errors.New("unexpected structure in balance sheet data map")
	}

	statement_data := make(map[string][]interface{})
	for field, fieldIndex := range statementDataMap {
		fieldIndexFloat, ok := fieldIndex.(float64)
		if !ok {
			return []map[string]interface{}{}, errors.New("unexpected index type in fieldIndex")
		}
		fieldValues := []interface{}{}
		for _, index := range data[int(fieldIndexFloat)].([]interface{}) {
			indexFloat, ok := index.(float64)
			if !ok {
				return []map[string]interface{}{}, errors.New("unexpected type in field values index")
			}
			fieldValues = append(fieldValues, data[int(indexFloat)])
		}
		statement_data[field] = fieldValues
	}

	// Converting the map of slices into a slice of maps to resemble final structure
	statement_data_slice := make([]map[string]interface{}, 0, len(statement_data["datekey"]))
	for i := 0; i < len(statement_data["datekey"]); i++ {
		record := make(map[string]interface{})
		for key, values := range statement_data {
			// There are some cases where the data is missing for a particular key
			if i >= len(values) {
				continue
			}
			record[key] = values[i]
		}
		statement_data_slice = append(statement_data_slice, record)
	}

	return statement_data_slice, nil
}

func GetBalanceSheets(symbol string) ([]BalanceSheet, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/financials/balance-sheet/__data.json", symbol)
	balanceSheetData, err := scrapeFinancialStatementData(url)
	if err != nil {
		return []BalanceSheet{}, err
	}

	balanceSheets := make([]BalanceSheet, 0, len(balanceSheetData))
	for i := 0; i < len(balanceSheetData); i++ {
		record := balanceSheetData[i]
		// Marshal the map to JSON
		jsonData, err := json.Marshal(record)
		if err != nil {
			return []BalanceSheet{}, err
		}
		// Unmarshal the JSON data into an instance of balanceSheet
		var balanceSheetRecord BalanceSheet
		err = json.Unmarshal(jsonData, &balanceSheetRecord)
		if err != nil {
			return []BalanceSheet{}, err
		}
		balanceSheets = append(balanceSheets, balanceSheetRecord)
	}
	return balanceSheets, nil
}

func GetCashFlows(symbol string) ([]CashFlow, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/financials/cash-flow-statement/__data.json?p=quarterly", symbol)
	cashFlowData, err := scrapeFinancialStatementData(url)
	if err != nil {
		return []CashFlow{}, err
	}

	cashFlows := make([]CashFlow, 0, len(cashFlowData))
	for i := 0; i < len(cashFlowData); i++ {
		record := cashFlowData[i]
		// Marshal the map to JSON
		jsonData, err := json.Marshal(record)
		if err != nil {
			return []CashFlow{}, err
		}
		// Unmarshal the JSON data into an instance of balanceSheet
		var cashFlowRecord CashFlow
		err = json.Unmarshal(jsonData, &cashFlowRecord)
		if err != nil {
			return []CashFlow{}, err
		}
		cashFlows = append(cashFlows, cashFlowRecord)
	}
	return cashFlows, nil
}

func GetIncomeStatements(symbol string) ([]IncomeStatement, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/financials/__data.json?p=quarterly", symbol)
	incomeStatementData, err := scrapeFinancialStatementData(url)
	if err != nil {
		return []IncomeStatement{}, err
	}

	incomeStatements := make([]IncomeStatement, 0, len(incomeStatementData))
	for i := 0; i < len(incomeStatementData); i++ {
		record := incomeStatementData[i]
		// Marshal the map to JSON
		jsonData, err := json.Marshal(record)
		if err != nil {
			return []IncomeStatement{}, err
		}
		// Unmarshal the JSON data into an instance of balanceSheet
		var incomeStatementRecord IncomeStatement
		err = json.Unmarshal(jsonData, &incomeStatementRecord)
		if err != nil {
			return []IncomeStatement{}, err
		}
		incomeStatements = append(incomeStatements, incomeStatementRecord)
	}
	return incomeStatements, nil
}

func GetFinancialRatios(symbol string) ([]FinancialRatios, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/financials/ratios/__data.json?p=quarterly", symbol)
	financialRatiosData, err := scrapeFinancialStatementData(url)
	if err != nil {
		return []FinancialRatios{}, err
	}

	financialRatios := make([]FinancialRatios, 0, len(financialRatiosData))
	for i := 0; i < len(financialRatiosData); i++ {
		record := financialRatiosData[i]
		// Marshal the map to JSON
		jsonData, err := json.Marshal(record)
		if err != nil {
			return []FinancialRatios{}, err
		}
		// Unmarshal the JSON data into an instance of balanceSheet
		var financialRatiosRecord FinancialRatios
		err = json.Unmarshal(jsonData, &financialRatiosRecord)
		if err != nil {
			return []FinancialRatios{}, err
		}
		financialRatios = append(financialRatios, financialRatiosRecord)
	}

	return financialRatios, nil
}
