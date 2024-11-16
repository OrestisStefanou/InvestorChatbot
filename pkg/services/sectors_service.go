package services

import (
	"fmt"
	"investbot/pkg/domain"
)

type SectorService interface {
	GenerateRagResponse(sessionId string, responseChannel chan<- string) error
}

type SectorDataService interface {
	GetSectorStocks(sector string) ([]domain.SectorStock, error)
	GetSectors() ([]domain.Sector, error)
}

type SessionService interface {
	GetConversationBySessionId(sessionId string) ([]map[string]string, error)
}

type SectorServiceRag struct {
	DataService    SectorDataService
	Llm            Llm
	SessionService SessionService
}

func (rag SectorServiceRag) GenerateRagResponse(sessionId string, responseChannel chan<- string) error {
	conversation, err := rag.SessionService.GetConversationBySessionId(sessionId)
	if err != nil {
		return err
	}
	fmt.Println(conversation)
	return nil
}
