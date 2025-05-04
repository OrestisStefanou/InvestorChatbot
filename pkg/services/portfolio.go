package services

import (
	"investbot/pkg/domain"
	"investbot/pkg/errors"
)

type PortfolioRepository interface {
	GetUserPortfolio(userEmail string) (domain.Portfolio, error)
	CreateUserPortfolio(portfolio domain.Portfolio) error
	UpdateUserPortfolio(portfolio domain.Portfolio) error
}

type PortfolioService struct {
	portfolioRepository PortfolioRepository
}

func NewPortfolioService(portfolioRepository PortfolioRepository) (*PortfolioService, error) {
	return &PortfolioService{portfolioRepository: portfolioRepository}, nil
}

func (s *PortfolioService) GetUserPortfolio(userEmail string) (domain.Portfolio, error) {
	portfolio, err := s.portfolioRepository.GetUserPortfolio(userEmail)
	if err != nil {
		return domain.Portfolio{}, err
	}

	if portfolio.UserEmail == "" {
		return domain.Portfolio{}, errors.PortfolioNotFoundError{UserEmail: userEmail}
	}

	return portfolio, nil
}

func (s *PortfolioService) CreateUserPortfolio(portfolio domain.Portfolio) error {
	return s.portfolioRepository.CreateUserPortfolio(portfolio)
}

func (s *PortfolioService) UpdateUserPortfolio(portfolio domain.Portfolio) error {
	portfolio, err := s.portfolioRepository.GetUserPortfolio(portfolio.UserEmail)
	if err != nil {
		return err
	}

	if portfolio.UserEmail == "" {
		return errors.PortfolioNotFoundError{UserEmail: portfolio.UserEmail}
	}

	return s.portfolioRepository.UpdateUserPortfolio(portfolio)
}
