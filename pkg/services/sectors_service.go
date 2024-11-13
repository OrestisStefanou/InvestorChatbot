package services

import "investbot/pkg/domain"

type SectorsService interface {
	GetAllSectors() ([]domain.Sector, error)
	GetSectorStocks(sector string) ([]domain.SectorStock, error)
	AskLLM(conversation []map[string]string, responseChannel chan<- string) error
}
