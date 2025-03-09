package main

import (
	"fmt"
	"investbot/pkg/marketDataScraper"
)

func main() {
	dataService := marketDataScraper.MarketDataScraper{}
	news, _ := dataService.GetStockNews("abnb")
	for _, v := range news {
		fmt.Printf("%+v\n", v)
	}
}
