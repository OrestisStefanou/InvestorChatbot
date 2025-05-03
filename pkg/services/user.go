package services

import (
	"investbot/pkg/domain"
)

type UserRepository interface {
	GetUser(id string) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) (*UserService, error) {
	return &UserService{userRepository: userRepository}, nil
}

func (s *UserService) GetUser(id string) (domain.User, error) {
	return s.userRepository.GetUser(id)
}

func (s *UserService) CreateUser(user domain.User) (domain.User, error) {
	return s.userRepository.CreateUser(user)
}
