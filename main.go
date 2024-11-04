package main

import (
	"fmt"
	"investbot/marketDataScraper"
)

func main() {
	etf, err := marketDataScraper.GetEtfOverview("spy")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", etf)
}
