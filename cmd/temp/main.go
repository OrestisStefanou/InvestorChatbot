package main

import (
	"fmt"
	"investbot/pkg/config"
	"investbot/pkg/marketDataScraper"
	"investbot/pkg/services"
)

type User struct {
	ID   int
	Name string
}

func main() {
	cache, _ := services.NewBadgerCacheService()
	conf, _ := config.LoadConfig()
	mds := marketDataScraper.NewMarketDataScraperWithCache(cache, conf)

	sectorStocks, err := mds.GetSectorStocks("technology")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sectorStocks[:10])

	fmt.Println("--------------------------------")

	sectorStocks, err = mds.GetSectorStocks("technology")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sectorStocks[:10])
}
