package main

import (
	"fmt"
	"investbot/pkg/marketDataScraper"
	"log"
)

func main() {
	scraper := marketDataScraper.MarketDataScraper{}
	portfolio, err := scraper.GetSuperInvestorPortfolio("Warren Buffett - Berkshire Hathaway")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("HOLDINGS")
	for _, v := range portfolio.Holdings {
		fmt.Printf("%+v\n", v)
	}

	fmt.Println("SECTOR ANALYSIS")
	for _, v := range portfolio.SectorAnalysis {
		fmt.Printf("%+v\n", v)
	}
}
