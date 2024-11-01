package main

import (
	"fmt"
	"investbot/marketDataScraper"
)

func main() {
	forecast, err := marketDataScraper.GetStockForecast("amzn")
	if err != nil {
		panic(err)
	}
	fmt.Println(forecast)
}
