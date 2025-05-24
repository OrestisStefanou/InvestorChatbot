package services

import (
	"investbot/pkg/domain"
	"investbot/pkg/errors"
	"time"
)

type PortfolioRepository interface {
	GetPortfolioById(portfolioID string) (domain.Portfolio, error)
	CreatePortfolio(portfolio domain.Portfolio) error
	UpdatePortfolio(portfolio domain.Portfolio) error
}

type PortfolioService struct {
	portfolioRepository PortfolioRepository
}

func NewPortfolioService(portfolioRepository PortfolioRepository) (*PortfolioService, error) {
	return &PortfolioService{portfolioRepository: portfolioRepository}, nil
}

func (s *PortfolioService) GetPortfolioById(portfolioID string) (domain.Portfolio, error) {
	portfolio, err := s.portfolioRepository.GetPortfolioById(portfolioID)
	if err != nil {
		return domain.Portfolio{}, err
	}

	if portfolio.ID == "" {
		return domain.Portfolio{}, errors.PortfolioNotFoundError{PortfolioID: portfolioID}
	}

	return portfolio, nil
}

func (s *PortfolioService) CreatePortfolio(portfolio domain.Portfolio) error {
	// TODO: Check that the symbols are valid and if duplicates exist
	portfolio.CreatedAt = time.Now()
	return s.portfolioRepository.CreatePortfolio(portfolio)
}

func (s *PortfolioService) UpdatePortfolio(portfolio domain.Portfolio) error {
	dbPortfolio, err := s.portfolioRepository.GetPortfolioById(portfolio.ID)
	if err != nil {
		return err
	}

	// TODO: Check that the symbols are valid and if duplicates exist
	portfolio.UpdatedAt = time.Now()
	portfolio.CreatedAt = dbPortfolio.CreatedAt
	return s.portfolioRepository.UpdatePortfolio(portfolio)
}
