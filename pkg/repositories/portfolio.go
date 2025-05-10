package repositories

import (
	"encoding/json"
	"errors"
	"investbot/pkg/domain"

	investbotErr "investbot/pkg/errors"

	badger "github.com/dgraph-io/badger/v4"
)

type PortfolioRepository struct {
	db *badger.DB
}

func NewPortfolioRepository(db *badger.DB) (*PortfolioRepository, error) {
	return &PortfolioRepository{db: db}, nil
}

func (r *PortfolioRepository) GetUserPortfolio(userEmail string) (domain.Portfolio, error) {
	var portfolio domain.Portfolio
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(userEmail))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return investbotErr.PortfolioNotFoundError{UserEmail: userEmail}
			}
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &portfolio)
		})
	})

	return portfolio, err
}

func (r *PortfolioRepository) CreateUserPortfolio(portfolio domain.Portfolio) error {
	err := r.db.Update(func(txn *badger.Txn) error {
		portfolioBytes, err := json.Marshal(portfolio)
		if err != nil {
			return err
		}

		return txn.Set([]byte(portfolio.UserEmail), portfolioBytes)
	})

	return err
}

func (r *PortfolioRepository) UpdateUserPortfolio(portfolio domain.Portfolio) error {
	err := r.db.Update(func(txn *badger.Txn) error {
		portfolioBytes, err := json.Marshal(portfolio)
		if err != nil {
			return err
		}

		return txn.Set([]byte(portfolio.UserEmail), portfolioBytes)
	})

	return err
}
