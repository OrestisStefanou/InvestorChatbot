package services

import (
	"errors"
	"investbot/pkg/domain"
	investbotErr "investbot/pkg/errors"
)

type UserContextRepository interface {
	GetUserContext(userID string) (domain.UserContext, error)
	InsertUserContext(userContext domain.UserContext) error
	UpdateUserContext(userContext domain.UserContext) error
}

type UserContextService struct {
	userContextRepository UserContextRepository
}

func NewUserContextService(userContextRepository UserContextRepository) (*UserContextService, error) {
	return &UserContextService{userContextRepository: userContextRepository}, nil
}

func (s *UserContextService) GetUserContext(userID string) (domain.UserContext, error) {
	userContext, err := s.userContextRepository.GetUserContext(userID)
	if err != nil {
		return domain.UserContext{}, err
	}

	return userContext, nil
}

func (s *UserContextService) CreateUserContext(userContext domain.UserContext) error {
	// Check if user context for given user id already exists
	dbUserContext, err := s.userContextRepository.GetUserContext(userContext.UserID)
	if err != nil {
		notFoundError := investbotErr.UserContextNotFoundError{UserID: dbUserContext.UserID}
		if errors.As(err, &notFoundError) {
			// do nothing in this case
		} else {
			return err
		}
	}

	if dbUserContext.UserID != "" {
		return investbotErr.UserContextAlreadyExistsError{UserID: dbUserContext.UserID}
	}

	// TODO: Check if symbols are valid and if duplicates exist
	return s.userContextRepository.InsertUserContext(userContext)
}

func (s *UserContextService) UpdateUserContext(userContext domain.UserContext) error {
	_, err := s.userContextRepository.GetUserContext(userContext.UserID)
	if err != nil { // user context not found error is covered here
		return err
	}
	return s.userContextRepository.UpdateUserContext(userContext)
}
