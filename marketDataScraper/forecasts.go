package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"investbot/domain"
	"io"
	"net/http"
)

func GetStockForecast(symbol string) (domain.StockForecast, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/forecast/__data.json", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return domain.StockForecast{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.StockForecast{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return domain.StockForecast{}, err
	}

	// Accessing nodes data
	nodes := rawData["nodes"].([]interface{})
	data := nodes[2].(map[string]interface{})["data"].([]interface{})
	dataMap := data[0].(map[string]interface{})

	// Estimates Scraping
	quarterlyEstimatesData := make(map[string][]interface{})

	estimatesDataIndex := int(dataMap["estimates"].(float64))
	estimatesTableDataIndex := int(data[estimatesDataIndex].(map[string]interface{})["table"].(float64))
	quarterlyEstimatesDataIndex := int(data[estimatesTableDataIndex].(map[string]interface{})["quarterly"].(float64))
	quarterlyEstimatesDataMap := data[quarterlyEstimatesDataIndex].(map[string]interface{})

	for estimationField, estimationFieldIdx := range quarterlyEstimatesDataMap {
		if estimationField == "lastDate" {
			continue
		}
		var estimationFieldValues []interface{}
		for _, fieldValueIndex := range data[int(estimationFieldIdx.(float64))].([]interface{}) {
			fieldValue := data[int(fieldValueIndex.(float64))]
			if fieldValue == "[PRO]" {
				continue
			}
			estimationFieldValues = append(estimationFieldValues, fieldValue)
		}
		quarterlyEstimatesData[estimationField] = estimationFieldValues
	}

	// Prepare estimations_doc in the same format as in Python
	var estimationsDoc []map[string]interface{}
	for i := 0; i < len(quarterlyEstimatesData["eps"]); i++ {
		record := make(map[string]interface{})
		for key, values := range quarterlyEstimatesData {
			record[key] = values[i]
		}
		estimationsDoc = append(estimationsDoc, record)
	}

	// Create a slice of StockEstimation structs from estimationsDoc
	estimations := make([]domain.StockEstimation, 0, len(estimationsDoc))
	for _, record := range estimationsDoc {
		// Check if the record keys are present
		var date string
		if record["dates"] == nil {
			date = ""
		} else {
			date = record["dates"].(string)
		}

		var eps float64
		if record["eps"] == nil {
			eps = 0
		} else {
			eps = record["eps"].(float64)
		}

		var epsGrowth float64
		if record["epsGrowth"] == nil {
			epsGrowth = 0
		} else {
			epsGrowth = record["epsGrowth"].(float64)
		}

		var revenue float64
		if record["revenue"] == nil {
			revenue = 0
		} else {
			revenue = record["revenue"].(float64)
		}

		var revenueGrowth float64
		if record["revenueGrowth"] == nil {
			revenueGrowth = 0
		} else {
			revenueGrowth = record["revenueGrowth"].(float64)
		}

		estimation := domain.StockEstimation{
			Date:          date,
			Eps:           eps,
			EpsGrowth:     epsGrowth,
			FiscalQuarter: record["fiscalQuarter"].(string),
			FiscalYear:    record["fiscalYear"].(string),
			Revenue:       revenue,
			RevenueGrowth: revenueGrowth,
		}
		estimations = append(estimations, estimation)
	}

	// Target Price Scraping
	targetDataIndex := int(dataMap["targets"].(float64))
	targetDataMap := data[targetDataIndex].(map[string]interface{})

	targetPriceKeys := []string{"average", "high", "low", "median"}
	targetPriceDoc := make(map[string]interface{})
	for _, targetKey := range targetPriceKeys {
		targetValueIndex := int(targetDataMap[targetKey].(float64))
		targetPriceDoc[targetKey] = data[targetValueIndex]
	}

	// Create a StockTargetPrc struct from targetPriceDoc
	targetPrice := domain.StockTargetPrc{
		Average: float32(targetPriceDoc["average"].(float64)),
		High:    float32(targetPriceDoc["high"].(float64)),
		Low:     float32(targetPriceDoc["low"].(float64)),
		Median:  float32(targetPriceDoc["median"].(float64)),
	}

	return domain.StockForecast{
		Estimations: estimations,
		TargetPrice: targetPrice,
	}, nil
}
