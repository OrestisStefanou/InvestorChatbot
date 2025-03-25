package services

import (
	"fmt"
	"investbot/pkg/errors"
	"investbot/pkg/services/faq"
	"math/rand"
)

type FaqTopic string

const (
	EDUCATION_FAQ_TOPIC      FaqTopic = "education"
	SECTORS_FAQ_TOPIC        FaqTopic = "sectors"
	STOCK_OVERVIEW_FAQ_TOPIC FaqTopic = "stock_overview"
	BALANCE_SHEET_FAQ_TOPIC  FaqTopic = "balance_sheet"
)

type FaqService struct {
	topicToFaq map[FaqTopic][]string
	faqLimit   int
}

func NewFaqService(faqLimit int) (*FaqService, error) {
	topicToFaq := map[FaqTopic][]string{
		EDUCATION_FAQ_TOPIC:      faq.EduationFaq,
		SECTORS_FAQ_TOPIC:        faq.SectorsFaq,
		STOCK_OVERVIEW_FAQ_TOPIC: faq.StockOverviewFaq,
		BALANCE_SHEET_FAQ_TOPIC:  faq.BalanceSheetFaq,
	}

	return &FaqService{topicToFaq: topicToFaq, faqLimit: faqLimit}, nil
}

func (s FaqService) GetFaqForTopic(topic FaqTopic) ([]string, error) {
	faqSlice, found := s.topicToFaq[topic]
	if !found {
		return nil, &errors.TopicNotFoundError{Message: fmt.Sprintf("Faq for %s not found", topic)}
	}

	randomFaq := s.pickRandomStrings(faqSlice, s.faqLimit)

	return randomFaq, nil
}

// pickRandomStrings selects k random strings from the given slice without modifying the original.
func (s FaqService) pickRandomStrings(original []string, k int) []string {
	n := len(original)
	if k > n {
		k = n // Ensure k does not exceed the length of original
	}

	// Create a new slice of length k to hold the result
	result := make([]string, k)

	// Shuffle only the necessary part
	for i := 0; i < k; i++ {
		j := rand.Intn(n-i) + i                             // Pick a random index from i to n-1
		original[i], original[j] = original[j], original[i] // Swap
		result[i] = original[i]                             // Copy the selected element into result
	}

	return result
}
