package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func streaming() {
	// Make a GET request to the URL
	resp, err := http.Get("https://postman-echo.com/stream/5")
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var chunk map[string]interface{} // or define a struct if you know the JSON structure
	for {
		if err := decoder.Decode(&chunk); err != nil {
			if err == io.EOF {
				break // end of stream
			}
			log.Fatalf("Failed to decode JSON chunk: %v", err)
		}
		fmt.Printf("Received JSON object: %v\n", chunk)
		fmt.Println("------------------------------")
	}
}
