package main

import (
	"investbot/pkg/config"
	"investbot/pkg/openAI"
	"investbot/pkg/services"
	"log"
)

func main() {
	config, _ := config.LoadConfig()
	openAiClient, _ := openAI.NewOpenAiClient(config.OpenAiKey, config.OpenAiBaseUrl)

	openAiLLM, err := openAI.NewOpenAiLLM(openAI.GPT4_MINI, openAiClient, 0.2)
	if err != nil {
		log.Fatal(err)
	}

	followUpQuestionsRag, _ := services.NewFollowUpQuestionsRag(openAiLLM)

	conversation := []services.Message{
		services.Message{
			Role:    services.User,
			Content: "Which are the top performing stock sectors?",
		},
		services.Message{
			Role: services.Assistant,
			Content: `The top-performing stock sectors can vary based on market conditions, economic cycles, and investor sentiment. However, as of 2024-2025, some of the best-performing stock sectors include:
1. Technology (Information Technology)
Why? AI, cloud computing, cybersecurity, and semiconductor growth.
Top Stocks: NVIDIA (NVDA), Microsoft (MSFT), Apple (AAPL), AMD (AMD), Broadcom (AVGO).
2. Energy (Oil & Renewables)
Why? Rising oil prices, energy transition to renewables.
Top Stocks: ExxonMobil (XOM), Chevron (CVX), NextEra Energy (NEE), Occidental Petroleum (OXY).
3. Healthcare (Biotech & Pharmaceuticals)
Why? Aging population, AI-driven drug discovery, weight-loss drugs (GLP-1 trend).
Top Stocks: Eli Lilly (LLY), Novo Nordisk (NVO), Pfizer (PFE), Moderna (MRNA).
4. Consumer Discretionary (Retail & E-commerce)
Why? Strong spending despite inflation, dominance of e-commerce.
Top Stocks: Amazon (AMZN), Tesla (TSLA), Home Depot (HD), McDonald's (MCD).
5. Industrials (Defense & Infrastructure)
Why? Global defense spending, U.S. infrastructure projects.
Top Stocks: Lockheed Martin (LMT), Raytheon (RTX), Caterpillar (CAT), Deere & Co. (DE).`,
		},
	}

	_, _ = followUpQuestionsRag.GenerateFollowUpQuestions(conversation, 5)
}
