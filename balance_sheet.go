package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	url := "https://stockanalysis.com/stocks/nvda/financials/balance-sheet/__data.json?p=quarterly"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Define the expected keys (not used here directly but may be useful in further filtering)
	expectedKeys := []string{"datekey", "fiscalYear", "fiscalQuarter", "cashneq", "investmentsc", "totalcash", "cashGrowth", "accountsReceivable", "otherReceivables", "receivables", "inventory", "restrictedCash", "othercurrent", "assetsc", "netPPE", "investmentsnc", "goodwill", "otherIntangibles", "othernoncurrent", "assets", "accountsPayable", "accruedExpenses", "debtc", "currentPortDebt", "currentCapLeases", "currentIncomeTaxesPayable", "currentUnearnedRevenue", "otherCurrentLiabilities", "currentLiabilities", "debtnc", "capitalLeases", "longTermUnearnedRevenue", "longTermDeferredTaxLiabilities", "otherliabilitiesnoncurrent", "liabilities", "commonStock", "retearn", "otherEquity", "equity", "liabilitiesequity", "sharesOutFilingDate", "sharesOutTotalCommon", "bvps", "tangibleBookValue", "tangibleBookValuePerShare", "debt", "netcash", "netCashGrowth", "netcashpershare", "workingcapital", "land", "machinery", "leaseholdImprovements", "tradingAssetSecurities"}

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		fmt.Println("Unexpected structure in 'nodes'")
		return
	}

	// Access the second element in "nodes"
	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		fmt.Println("Unexpected structure in 'nodes[2]'")
		return
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok {
		fmt.Println("Unexpected structure in 'data'")
		return
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		fmt.Println("Unexpected structure in 'data[0]'")
		return
	}

	financialDataIndex, ok := dataMap["financialData"].(float64)
	if !ok {
		fmt.Println("Unexpected structure for 'financialData'")
		return
	}

	// Retrieve data map at financial data index
	balanceSheetDataMap, ok := data[int(financialDataIndex)].(map[string]interface{})
	if !ok {
		fmt.Println("Unexpected structure in balance sheet data map")
		return
	}

	balanceSheetData := make(map[string][]interface{})
	for field, fieldIndex := range balanceSheetDataMap {
		fieldIndexFloat, ok := fieldIndex.(float64)
		if !ok {
			fmt.Println("Unexpected index type in fieldIndex")
			return
		}
		fieldValues := []interface{}{}
		for _, index := range data[int(fieldIndexFloat)].([]interface{}) {
			indexFloat, ok := index.(float64)
			if !ok {
				fmt.Println("Unexpected type in field values index")
				return
			}
			fieldValues = append(fieldValues, data[int(indexFloat)])
		}
		balanceSheetData[field] = fieldValues
	}

	// Converting the map into a slice of maps to resemble final structure
	result := []map[string]interface{}{}
	for i := 0; i < len(balanceSheetData[expectedKeys[0]]); i++ {
		record := make(map[string]interface{})
		for key, values := range balanceSheetData {
			record[key] = values[i]
		}
		result = append(result, record)
	}

	// Output final data
	fmt.Println(len(result))
	fmt.Println("--------------------------")
	fmt.Println(result[0]["totalcash"])
}
