package main

import (
	"fmt"
	"investbot/marketDataScraper"
)

func main() {
	balance_sheets, err := marketDataScraper.GetFinancialRatios("nvda")
	if err != nil {
		panic(err)

	}

	fmt.Println(balance_sheets)
}
