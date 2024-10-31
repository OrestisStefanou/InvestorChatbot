package main

import (
	"encoding/json"
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

func get_financial_statement_data() []map[string]interface{} {
	statement_data := make(map[string][]interface{})
	statement_data_slice := []map[string]interface{}{}
	url := "https://stockanalysis.com/stocks/aapl/financials/balance-sheet/__data.json?p=quarterly"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return statement_data_slice
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return statement_data_slice
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return statement_data_slice
	}

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		fmt.Println("Unexpected structure in 'nodes'")
		return statement_data_slice
	}

	// Access the second element in "nodes" which contains the data we are interested in
	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		fmt.Println("Unexpected structure in 'nodes[2]'")
		return statement_data_slice
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok {
		fmt.Println("Unexpected structure in 'data'")
		return statement_data_slice
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		fmt.Println("Unexpected structure in 'data[0]'")
		return statement_data_slice
	}

	financialDataIndex, ok := dataMap["financialData"].(float64)
	if !ok {
		fmt.Println("Unexpected structure for 'financialData'")
		return statement_data_slice
	}

	// Retrieve data map at financial data index
	statementDataMap, ok := data[int(financialDataIndex)].(map[string]interface{})
	if !ok {
		fmt.Println("Unexpected structure in balance sheet data map")
		return statement_data_slice
	}

	for field, fieldIndex := range statementDataMap {
		fieldIndexFloat, ok := fieldIndex.(float64)
		if !ok {
			fmt.Println("Unexpected index type in fieldIndex")
			return statement_data_slice
		}
		fieldValues := []interface{}{}
		for _, index := range data[int(fieldIndexFloat)].([]interface{}) {
			indexFloat, ok := index.(float64)
			if !ok {
				fmt.Println("Unexpected type in field values index")
				return statement_data_slice
			}
			fieldValues = append(fieldValues, data[int(indexFloat)])
		}
		statement_data[field] = fieldValues
	}

	// Converting the map of slices into a slice of maps to resemble final structure
	for i := 0; i < len(statement_data["datekey"]); i++ {
		record := make(map[string]interface{})
		for key, values := range statement_data {
			record[key] = values[i]
		}
		statement_data_slice = append(statement_data_slice, record)
	}
	fmt.Println(statement_data_slice)
	return statement_data_slice
}

func get_balance_sheets() {
	// url := "https://stockanalysis.com/stocks/aapl/financials/balance-sheet/__data.json?p=quarterly"
	// resp, err := http.Get(url)
	// if err != nil {
	// 	fmt.Println("Error fetching data:", err)
	// 	return
	// }
	// defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Error reading response body:", err)
	// 	return
	// }

	// var rawData map[string]interface{}
	// if err := json.Unmarshal(body, &rawData); err != nil {
	// 	fmt.Println("Error unmarshalling JSON:", err)
	// 	return
	// }

	// // Extract "nodes" from rawData
	// nodes, ok := rawData["nodes"].([]interface{})
	// if !ok || len(nodes) < 3 {
	// 	fmt.Println("Unexpected structure in 'nodes'")
	// 	return
	// }

	// // Access the second element in "nodes" which contains the data we are interested in
	// nodeData, ok := nodes[2].(map[string]interface{})
	// if !ok {
	// 	fmt.Println("Unexpected structure in 'nodes[2]'")
	// 	return
	// }

	// data, ok := nodeData["data"].([]interface{})
	// if !ok {
	// 	fmt.Println("Unexpected structure in 'data'")
	// 	return
	// }

	// dataMap, ok := data[0].(map[string]interface{})
	// if !ok {
	// 	fmt.Println("Unexpected structure in 'data[0]'")
	// 	return
	// }

	// financialDataIndex, ok := dataMap["financialData"].(float64)
	// if !ok {
	// 	fmt.Println("Unexpected structure for 'financialData'")
	// 	return
	// }

	// // Retrieve data map at financial data index
	// balanceSheetDataMap, ok := data[int(financialDataIndex)].(map[string]interface{})
	// if !ok {
	// 	fmt.Println("Unexpected structure in balance sheet data map")
	// 	return
	// }

	// balanceSheetData := make(map[string][]interface{})
	// for field, fieldIndex := range balanceSheetDataMap {
	// 	fieldIndexFloat, ok := fieldIndex.(float64)
	// 	if !ok {
	// 		fmt.Println("Unexpected index type in fieldIndex")
	// 		return
	// 	}
	// 	fieldValues := []interface{}{}
	// 	for _, index := range data[int(fieldIndexFloat)].([]interface{}) {
	// 		indexFloat, ok := index.(float64)
	// 		if !ok {
	// 			fmt.Println("Unexpected type in field values index")
	// 			return
	// 		}
	// 		fieldValues = append(fieldValues, data[int(indexFloat)])
	// 	}
	// 	balanceSheetData[field] = fieldValues
	// }

	balanceSheetData := get_financial_statement_data()
	for i := 0; i < len(balanceSheetData); i++ {
		record := balanceSheetData[i]
		// Marshal the map to JSON
		jsonData, err := json.Marshal(record)
		if err != nil {
			fmt.Println("Error marshaling data:", err)
			return
		}
		// Unmarshal the JSON data into an instance of balanceSheet
		var balanceSheetRecord BalanceSheet
		err = json.Unmarshal(jsonData, &balanceSheetRecord)
		if err != nil {
			fmt.Println("Error unmarshaling data:", err)
			return
		}
		fmt.Println(balanceSheetRecord)
	}
}

func main() {
	get_balance_sheets()
}
