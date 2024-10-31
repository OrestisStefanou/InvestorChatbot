package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Industry struct {
	name             string
	urlName          string
	numberOfStocks   int
	marketCap        float32
	dividendYieldPct float32
	peRatio          float32
	profitMarginPct  float32
	oneYearChangePct float32
}

func get_industries() {
	url := "https://stockanalysis.com/stocks/industry/all/__data.json"
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
	industriesDataIndex, ok := dataMap["industries"].(float64)
	if !ok {
		fmt.Println("Unexpected structure in 'data[0]'")
		fmt.Println(err)
		return
	}

	industryDataIndicesArray := data[int(industriesDataIndex)].([]interface{})

	industries := make([]Industry, 0, len(industryDataIndicesArray))
	for i := 0; i < len(industryDataIndicesArray); i++ {
		industryDataIndex := int(industryDataIndicesArray[i].(float64))
		industryData := data[industryDataIndex].(map[string]interface{})
		industryNameIndex := int(industryData["industry_name"].(float64))
		industryUrlNameIndex := int(industryData["url"].(float64))
		numberOfStocksIndex := int(industryData["stocks"].(float64))
		marketCapIndex := int(industryData["marketCap"].(float64))
		profitMarginIndex := int(industryData["profitMargin"].(float64))
		oneYearChangeIndex := int(industryData["ch1y"].(float64))

		// peRatio and dividendYield are handled differently because they could be missing from the industryData map
		var peRatio float32
		peRatioIndex, ok := industryData["peRatio"]
		if !ok {
			peRatio = 0
		} else {
			peRatioIndexInt := int(peRatioIndex.(float64))
			peRatio = float32(data[peRatioIndexInt].(float64))
		}

		var dividendYield float32
		dividendYieldIndex, ok := industryData["dividendYield"]
		if !ok {
			dividendYield = 0
		} else {
			dividendYieldIndexInt := int(dividendYieldIndex.(float64))
			dividendYield = float32(data[dividendYieldIndexInt].(float64))
		}

		industry := Industry{
			name:             data[industryNameIndex].(string),
			urlName:          data[industryUrlNameIndex].(string),
			numberOfStocks:   int(data[numberOfStocksIndex].(float64)),
			marketCap:        float32(data[marketCapIndex].(float64)),
			dividendYieldPct: dividendYield,
			peRatio:          peRatio,
			profitMarginPct:  float32(data[profitMarginIndex].(float64)),
			oneYearChangePct: float32(data[oneYearChangeIndex].(float64)),
		}
		industries = append(industries, industry)
	}

	fmt.Println(industries)
}

func main() {
	get_industries()
}
