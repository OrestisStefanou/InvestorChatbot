package repositories

import (
	"encoding/json"
	"investbot/pkg/domain"

	badger "github.com/dgraph-io/badger/v4"
)

type UserRepository struct {
	db *badger.DB
}

func NewUserRepository(db *badger.DB) (*UserRepository, error) {
	return &UserRepository{db: db}, nil
}

func (r *UserRepository) GetUser(email string) (domain.User, error) {
	var user domain.User
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(email))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})
	})

	return user, err
}

func (r *UserRepository) CreateUser(user domain.User) (domain.User, error) {
	err := r.db.Update(func(txn *badger.Txn) error {
		userBytes, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return txn.Set([]byte(user.Email), userBytes)
	})

	return user, err
}
