package errors

import "fmt"

type PortfolioNotFoundError struct {
	UserEmail string
}

func (e PortfolioNotFoundError) Error() string {
	return fmt.Sprintf("portfolio not found for user %s", e.UserEmail)
}
