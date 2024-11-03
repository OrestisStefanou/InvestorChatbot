package main

import (
	"fmt"
	"investbot/marketDataScraper"
)

func main() {
	etfs, err := marketDataScraper.GetEtfs()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(etfs[0])
}
