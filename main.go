package main

import (
	"fmt"
	"investbot/marketDataScraper"
)

func main() {
	etf, err := marketDataScraper.GetIndustryStocks("home-improvement-retail")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", etf)
}
