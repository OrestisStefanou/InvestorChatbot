package errors

import "fmt"

type PortfolioNotFoundError struct {
	PortfolioID string
}

func (e PortfolioNotFoundError) Error() string {
	return fmt.Sprintf("portfolio not found for id %s", e.PortfolioID)
}
