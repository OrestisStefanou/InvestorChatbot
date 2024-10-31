package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Sector struct {
	name             string
	urlName          string
	numberOfStocks   int
	marketCap        float32
	dividendYieldPct float32
	peRatio          float32
	profitMarginPct  float32
	oneYearChangePct float32
}

func get_sectors() {
	url := "https://stockanalysis.com/stocks/industry/sectors/__data.json"
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

	fmt.Println(dataMap)
	sectorsDataIndex, ok := dataMap["sectors"].(float64)
	if !ok {
		fmt.Println("Unexpected structure in 'data[0]'")
		fmt.Println(err)
		return
	}

	sectorDataIndicesArray := data[int(sectorsDataIndex)].([]interface{})

	sectors := make([]Sector, 0, len(sectorDataIndicesArray))
	for i := 0; i < len(sectorDataIndicesArray); i++ {
		sectorDataIndex := int(sectorDataIndicesArray[i].(float64))
		sectorData := data[sectorDataIndex].(map[string]interface{})
		sectorNameIndex := int(sectorData["sector_name"].(float64))
		sectorUrlNameIndex := int(sectorData["url"].(float64))
		numberOfStocksIndex := int(sectorData["stocks"].(float64))
		marketCapIndex := int(sectorData["marketCap"].(float64))
		dividendYieldIndex := int(sectorData["dividendYield"].(float64))
		peRatioIndex := int(sectorData["peRatio"].(float64))
		profitMarginIndex := int(sectorData["profitMargin"].(float64))
		oneYearChangeIndex := int(sectorData["ch1y"].(float64))

		sector := Sector{
			name:             data[sectorNameIndex].(string),
			urlName:          data[sectorUrlNameIndex].(string),
			numberOfStocks:   int(data[numberOfStocksIndex].(float64)),
			marketCap:        float32(data[marketCapIndex].(float64)),
			dividendYieldPct: float32(data[dividendYieldIndex].(float64)),
			peRatio:          float32(data[peRatioIndex].(float64)),
			profitMarginPct:  float32(data[profitMarginIndex].(float64)),
			oneYearChangePct: float32(data[oneYearChangeIndex].(float64)),
		}
		sectors = append(sectors, sector)
	}

	fmt.Println(sectors)
}
