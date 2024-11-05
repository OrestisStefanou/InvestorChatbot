package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"investbot/domain"
	"io"
	"net/http"
)

func GetStockProfile(symbol string) (domain.StockProfile, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/company/__data.json", symbol)

	resp, err := http.Get(url)
	if err != nil {
		return domain.StockProfile{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.StockProfile{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return domain.StockProfile{}, err
	}

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure in 'nodes'")
	}

	// Access the second element in "nodes" which contains the data we are interested in
	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure in 'nodes[2]'")
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure in 'data'")
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure in 'data[0]'")
	}

	descriptionIndex, ok := dataMap["description"].(float64)
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure for 'description'")
	}

	profileIndex, ok := dataMap["profile"].(float64)
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure for 'profile'")
	}

	stockProfileData := data[int(profileIndex)].(map[string]interface{})

	industryDataIndex, ok := stockProfileData["industry"].(float64)
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure for 'industry'")
	}

	industryData := data[int(industryDataIndex)].(map[string]interface{})
	industryNameIndex := industryData["value"].(float64)

	sectorDataIndex, ok := stockProfileData["sector"].(float64)
	if !ok {
		return domain.StockProfile{}, fmt.Errorf("unexpected structure for 'sector'")
	}
	sectorData := data[int(sectorDataIndex)].(map[string]interface{})
	sectorNameIndex := sectorData["value"].(float64)

	stockNameINdex := stockProfileData["name"].(float64)
	stockCountryIndex := stockProfileData["country"].(float64)
	stockFoundedIndex := stockProfileData["founded"].(float64)
	stockIpoDateIndex := stockProfileData["ipoDate"].(float64)
	stockCeoIndex := stockProfileData["ceo"].(float64)

	return domain.StockProfile{
		Name:        data[int(stockNameINdex)].(string),
		Description: data[int(descriptionIndex)].(string),
		Country:     data[int(stockCountryIndex)].(string),
		Founded:     int(data[int(stockFoundedIndex)].(float64)),
		IpoDate:     data[int(stockIpoDateIndex)].(string),
		Industry:    data[int(industryNameIndex)].(string),
		Sector:      data[int(sectorNameIndex)].(string),
		Ceo:         data[int(stockCeoIndex)].(string),
	}, nil
}
