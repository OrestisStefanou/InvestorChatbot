package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func scrapeSuperInvestorsAndPortfolioLinks() {
	url := "https://www.dataroma.com/m/managers.php"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	rsp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch the page: %v", err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		log.Fatalf("Error: status code %d", rsp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		log.Fatalf("Failed to parse the page: %v", err)
	}

	// Target the specific table containing super investors
	doc.Find("table#grid a[href^='/m/holdings.php']").Each(func(index int, item *goquery.Selection) {
		name := item.Text()
		link, exists := item.Attr("href")
		if exists && len(name) > 0 {
			fmt.Printf("Investor: %s, Link: %s\n", name, "https://www.dataroma.com"+link)
		}
	})
}

func main() {
	// Make GET request
	url := "https://www.dataroma.com/m/holdings.php?m=GFT"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Error: status code %d", resp.StatusCode)
	}

	// Parse HTML response with goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error parsing HTML:", err)
	}

	// // Find all tables in the document
	// doc.Find("table").Each(func(index int, tableHtml *goquery.Selection) {
	// 	// Print each table (you can customize this part to extract specific data)
	// 	fmt.Println("Table", index+1, ":")
	// 	tableHtml.Each(func(rowIndex int, rowHtml *goquery.Selection) {
	// 		rowHtml.Find("td").Each(func(cellIndex int, cellHtml *goquery.Selection) {
	// 			fmt.Print(cellHtml.Text(), "\t")
	// 		})
	// 		fmt.Println()
	// 	})
	// })
	// Find all tables in the document
	doc.Find("table").Each(func(index int, tableHtml *goquery.Selection) {
		// For each table, extract rows and print in JSON format
		tableHtml.Find("tr").Each(func(rowIndex int, rowHtml *goquery.Selection) {
			// Skip empty rows (optional)
			if rowHtml.Children().Length() == 0 {
				return
			}

			// Create a map to hold the row data
			rowData := make(map[string]interface{})

			// For each cell in the row, map column to value (you can customize this part based on your table structure)
			rowHtml.Find("td").Each(func(cellIndex int, cellHtml *goquery.Selection) {
				// You can customize the keys based on the columns' index or content (e.g., Column 1, Column 2)
				columnKey := fmt.Sprintf("Column%d", cellIndex+1)
				cellValue := cellHtml.Text()
				rowData[columnKey] = cellValue
			})

			// Marshal the map into JSON
			rowJSON, err := json.Marshal(rowData)
			if err != nil {
				log.Println("Error marshaling row data to JSON:", err)
				return
			}

			// Print the JSON string
			fmt.Println(string(rowJSON))
		})
	})
}
