package main

import (
	"fmt"
	"investbot/pkg/config"
	"investbot/pkg/openAI"
	"investbot/pkg/services/prompts"
	"log"
)

func main() {
	config, _ := config.LoadConfig()
	openAiClient, _ := openAI.NewOpenAiClient(config.OpenAiKey, "https://api.openai.com/v1")

	openAiLLM := openAI.OpenAiLLM{
		ModelName:   "gpt-4o-mini",
		Client:      openAiClient,
		Temperature: 0.2,
	}

	//dataService := marketDataScraper.MarketDataScraper{}
	//sessionService := services.MockSessionService{}
	//sectorService := services.EducationServiceRag{
	// DataService:    dataService,
	//Llm:            openAiLLM,
	//SessionService: sessionService,
	//}
	chunkChannel := make(chan string)

	go func() {
		conversation := []map[string]string{
			{"role": "system", "content": prompts.EducationPrompt},
			{"role": "user", "content": "Should I invest in stocks or ETFs?"},
			{"role": "system", "content": `Deciding between stocks and ETFs depends on your investment goals and preferences. Here’s a simple breakdown:
				**Stocks:**
				- **Ownership:** When you buy a stock, you own a piece of a company.
				- **Potential for High Returns:** Individual stocks can offer high returns if the company performs well.
				- **Higher Risk:** Investing in individual stocks can be riskier because their prices can be volatile.

				**ETFs (Exchange-Traded Funds):**
				- **Diversification:** ETFs hold a collection of stocks or other assets, which spreads out your risk.
				- **Lower Risk:** Because they are diversified, ETFs can be less risky than individual stocks.
				- **Ease of Trading:** ETFs trade like stocks on an exchange, making them easy to buy and sell.

				**Considerations:**
				- If you prefer to pick individual companies and are comfortable with higher risk, stocks might be for you.
				- If you want a more diversified investment with potentially lower risk, ETFs could be a better choice.

				Ultimately, it might be beneficial to include both in your portfolio based on your risk tolerance and investment strategy!!`},
			{"role": "user", "content": "What about crypto?"},
		}
		if err := openAiLLM.GenerateResponse(conversation, chunkChannel); err != nil {
			// Handle the error (e.g., log it)
			log.Printf("Error during request: %v", err)
			close(chunkChannel) // Ensure the channel is closed if there’s an error
		}
	}()

	// Consume the chunks from the channel
	var finalResponse string
	for content := range chunkChannel {
		// Process the chunk as it arrives
		fmt.Printf("Received chunk: %s\n", content)
		finalResponse += content
	}

	// Optional: After the channel is closed, perform any final tasks
	log.Println("Streaming has finished.")
	fmt.Println(finalResponse)

}
