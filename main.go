package main

import (
	"fmt"
	"investbot/marketDataScraper"
)

func main() {
	etf, err := marketDataScraper.GetBalanceSheets("aapl")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", etf)
}
