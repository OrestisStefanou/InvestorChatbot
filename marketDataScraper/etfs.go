package marketDataScraper

import (
	"encoding/json"
	"net/http"
)

type Etf struct {
	Symbol     string
	Name       string
	AssetClass string
	Aum        float32
}

func GetEtfs() ([]Etf, error) {
	url := "https://api.stockanalysis.com/api/screener/e/f?m=s&s=asc&c=s,n,assetClass,aum&i=etf"

	resp, err := http.Get(url)
	if err != nil {
		return []Etf{}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Etf{}, nil
	}

	// Define an anonymous struct to match the JSON structure
	var apiResponse struct {
		Status int `json:"status"`
		Data   struct {
			Data []struct {
				S          string  `json:"s"`
				N          string  `json:"n"`
				AssetClass string  `json:"assetClass"`
				AUM        float64 `json:"aum"`
			} `json:"data"`
			ResultsCount int `json:"resultsCount"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return []Etf{}, nil
	}

	etfs := make([]Etf, 0, len(apiResponse.Data.Data))
	for _, etfData := range apiResponse.Data.Data {
		etf := Etf{
			Symbol:     etfData.S,
			Name:       etfData.N,
			AssetClass: etfData.AssetClass,
			Aum:        float32(etfData.AUM),
		}
		etfs = append(etfs, etf)
	}
	return etfs, nil
}
