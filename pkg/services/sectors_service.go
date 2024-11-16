package services

import (
	"fmt"
	"investbot/pkg/domain"
)

type SectorService interface {
	GenerateRagResponse(sessionId string, question string, responseChannel chan<- string) error
}

type SectorDataService interface {
	GetSectorStocks(sector string) ([]domain.SectorStock, error)
	GetSectors() ([]domain.Sector, error)
}

type sectorContext struct {
	sector       domain.Sector
	sectorStocks []domain.SectorStock
}

type SectorServiceRag struct {
	DataService    SectorDataService
	Llm            Llm
	SessionService SessionService
	Prompt         string
}

func (rag SectorServiceRag) createRagContext() (string, error) {
	var ragContext string
	sectors, err := rag.DataService.GetSectors()
	if err != nil {
		return ragContext, err
	}

	for i := 0; i < len(sectors); i++ {
		sector := sectors[i]
		sectorStocks, err := rag.DataService.GetSectorStocks(sector.UrlName)
		if err != nil {
			return ragContext, err
		}
		context := sectorContext{
			sector:       sector,
			sectorStocks: sectorStocks,
		}
		ragContext += fmt.Sprintf("%+v\n", context)
	}

	return ragContext, nil
}

func (rag SectorServiceRag) GenerateRagResponse(sessionId string, question string, responseChannel chan<- string) error {
	conversation, err := rag.SessionService.GetConversationBySessionId(sessionId)
	if err != nil {
		return err
	}

	if len(conversation) == 0 {
		// Format the prompt to contain the neccessary context
		ragContext, err := rag.createRagContext()
		if err != nil {
			return err
		}
		prompt := fmt.Sprintf(rag.Prompt, ragContext)
		conversation = append(conversation, map[string]string{"role": "system", "content": prompt})
	}
	conversation = append(conversation, map[string]string{"role": "user", "content": question})

	if err := rag.Llm.GenerateResponse(conversation, responseChannel); err != nil {
		return err
	}

	return nil
}
