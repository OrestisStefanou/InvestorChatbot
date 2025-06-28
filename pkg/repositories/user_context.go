package repositories

import (
	"encoding/json"
	"errors"
	"investbot/pkg/domain"
	investbotErr "investbot/pkg/errors"

	"github.com/dgraph-io/badger/v4"
)

type UserContextRepository struct {
	db *badger.DB
}

func NewUserContextRepository(db *badger.DB) (*UserContextRepository, error) {
	return &UserContextRepository{db: db}, nil
}

func (r *UserContextRepository) GetUserContext(userID string) (domain.UserContext, error) {
	var userContext domain.UserContext
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(userID))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return investbotErr.UserContextNotFoundError{UserID: userID}
			}
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &userContext)
		})
	})

	return userContext, err
}

func (r *UserContextRepository) InsertUserContext(userContext domain.UserContext) error {
	err := r.db.Update(func(txn *badger.Txn) error {
		userContextBytes, err := json.Marshal(userContext)
		if err != nil {
			return err
		}

		return txn.Set([]byte(userContext.UserID), userContextBytes)
	})

	return err
}

func (r *UserContextRepository) UpdateUserContext(userContext domain.UserContext) error {
	err := r.db.Update(func(txn *badger.Txn) error {
		userContextBytes, err := json.Marshal(userContext)
		if err != nil {
			return err
		}

		return txn.Set([]byte(userContext.UserID), userContextBytes)
	})

	return err
}
