package services

import (
	"investbot/pkg/domain"
)

type UserService interface {
	GetUser(id string) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
}
