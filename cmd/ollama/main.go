package main

import (
	"fmt"
	"investbot/pkg/marketDataScraper"
)

func main() {
	dataService := marketDataScraper.MarketDataScraper{}
	forecast, _ := dataService.GetStockForecast("nvda")
	fmt.Printf("Target Price: %+v\n", forecast.TargetPrice)
	fmt.Println("ESTIMATIONS")
	for _, v := range forecast.Estimations {
		fmt.Printf("%+v\n", v)
	}
}
