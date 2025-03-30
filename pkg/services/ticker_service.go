package services

import "investbot/pkg/domain"

type TickerDataService interface {
	GetTickers() ([]domain.Ticker, error)
}

type TickerService struct {
	dataService TickerDataService
}

func NewTickerServiceImpl(dataService TickerDataService) (*TickerService, error) {
	return &TickerService{
		dataService: dataService,
	}, nil
}

type TickerFilterOptions struct {
	Limit        int
	Page         int
	SearchString string
}

func (s TickerService) GetTickers(filters TickerFilterOptions) ([]domain.Ticker, error) {
	// TODO: Implement the filtering
	return s.dataService.GetTickers()
}
