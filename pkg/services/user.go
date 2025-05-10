package services

import (
	"investbot/pkg/domain"
	"investbot/pkg/errors"
)

type UserRepository interface {
	GetUser(email string) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) (*UserService, error) {
	return &UserService{userRepository: userRepository}, nil
}

func (s *UserService) GetUser(email string) (domain.User, error) {
	return s.userRepository.GetUser(email)
}

func (s *UserService) CreateUser(user domain.User) (domain.User, error) {
	existingUser, err := s.GetUser(user.Email)
	if err != nil {
		return domain.User{}, err
	}

	if existingUser.Email != "" {
		return domain.User{}, errors.UserAlreadyExistsError{Message: "User with this email already exists"}
	}
	return s.userRepository.CreateUser(user)
}
