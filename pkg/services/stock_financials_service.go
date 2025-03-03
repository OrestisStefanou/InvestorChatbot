package services

import (
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
	"time"
)

type StockFinancialsDataService interface {
	GetBalanceSheets(symbol string) ([]domain.BalanceSheet, error)
	GetIncomeStatements(symbol string) ([]domain.IncomeStatement, error)
	GetCashFlows(symbol string) ([]domain.CashFlow, error)
}

type stockFinancialsContext struct {
	currentDate      string
	symbol           string
	balanceSheets    []domain.BalanceSheet
	incomeStatements []domain.IncomeStatement
	cashFlows        []domain.CashFlow
}

type StockFinancialsRag struct {
	dataService StockFinancialsDataService
	llm         Llm
}

func NewStockFinancialsRag(llm Llm, stockFinancialsDataService StockFinancialsDataService) (*StockFinancialsRag, error) {
	return &StockFinancialsRag{dataService: stockFinancialsDataService, llm: llm}, nil
}

func (rag StockFinancialsRag) createRagContext(tags Tags) (string, error) {
	var ragContext stockFinancialsContext
	ragContext.symbol = tags.StockSymbol
	ragContext.currentDate = time.Now().Format("2006-01-02")

	if tags.BalanceSheet {
		balanceSheets, err := rag.dataService.GetBalanceSheets(tags.StockSymbol)
		if err != nil {
			return "", &DataServiceError{Message: fmt.Sprintf("GetBalanceSheets failed: %s", err)}
		}
		ragContext.balanceSheets = balanceSheets
	}

	if tags.CashFlow {
		cashFlows, err := rag.dataService.GetCashFlows(tags.StockSymbol)
		if err != nil {
			return "", &DataServiceError{Message: fmt.Sprintf("GetCashFlows failed: %s", err)}
		}
		ragContext.cashFlows = cashFlows
	}

	if tags.IncomeStatement {
		incomeStatements, err := rag.dataService.GetIncomeStatements(tags.StockSymbol)
		if err != nil {
			return "", &DataServiceError{Message: fmt.Sprintf("GetIncomeStatements failed: %s", err)}
		}
		ragContext.incomeStatements = incomeStatements
	}

	return fmt.Sprintf("%+v\n", ragContext), nil
}

func (rag StockFinancialsRag) GenerateRagResponse(conversation []Message, tags Tags, responseChannel chan<- string) error {
	// Format the prompt to contain the neccessary context
	ragContext, err := rag.createRagContext(tags)
	if err != nil {
		return err
	}

	prompt := fmt.Sprintf(prompts.StockFinancialsPrompt, ragContext)
	prompt_msg := Message{
		Role:    System,
		Content: prompt,
	}

	// Add the prompt as the first message in the existing conversation
	conversation_with_prompt := append([]Message{prompt_msg}, conversation...)

	if err := rag.llm.GenerateResponse(conversation_with_prompt, responseChannel); err != nil {
		return &RagError{Message: fmt.Sprintf("GenerateResponse failed: %s", err)}
	}

	return nil
}
