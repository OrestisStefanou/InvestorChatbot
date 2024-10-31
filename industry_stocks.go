package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IndustryStock struct {
	symbol      string
	companyName string
	marketCap   float32
}

func get_industry_stocks() {
	url := "https://stockanalysis.com/stocks/industry/biotechnology/__data.json?"
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

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		fmt.Println("Unexpected structure in 'nodes'")
		return
	}

	// Access the second element in "nodes" which contains the data we are interested in
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

	stocksArrayDataIndex, ok := dataMap["data"].(float64)
	if !ok {
		fmt.Println("Unexpected structure in 'data[0]'")
		return
	}

	stocksDataIndicesArray := data[int(stocksArrayDataIndex)].([]interface{})
	if !ok {
		fmt.Println("Unexpected structure in 'data[0]'")
		return
	}

	stocks := make([]IndustryStock, len(stocksDataIndicesArray))
	for i := 0; i < len(stocksDataIndicesArray); i++ {
		stockDataIndex := int(stocksDataIndicesArray[i].(float64))
		stockData := data[stockDataIndex].(map[string]interface{})
		stockSymbolIndex := int(stockData["s"].(float64))
		stockCompanyNameIndex := int(stockData["n"].(float64))
		stockMarketCapIndex := int(stockData["marketCap"].(float64))
		stock := IndustryStock{
			symbol:      data[stockSymbolIndex].(string),
			companyName: data[stockCompanyNameIndex].(string),
			marketCap:   float32(data[stockMarketCapIndex].(float64)),
		}
		stocks = append(stocks, stock)
	}
	fmt.Println(stocks)
}

func main() {
	get_industry_stocks()
}
