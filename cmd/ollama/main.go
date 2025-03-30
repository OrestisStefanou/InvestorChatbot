package main

import (
	"fmt"
	"investbot/pkg/marketDataScraper"
	"log"
)

func main() {
	scraper := marketDataScraper.MarketDataScraper{}
	tickers, err := scraper.GetTickers()

	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, v := range tickers {
		fmt.Println(v)
	}
}
