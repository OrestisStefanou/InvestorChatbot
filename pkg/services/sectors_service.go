package services

import (
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
)

type SectorDataService interface {
	GetSectorStocks(sector string) ([]domain.SectorStock, error)
	GetSectors() ([]domain.Sector, error)
}

type sectorContext struct {
	sector       domain.Sector
	sectorStocks []domain.SectorStock
}

type SectorRag struct {
	BaseRag
	dataService        SectorDataService
	userContextService UserContextDataService
}

func NewSectorRag(
	llm Llm,
	sectorDataService SectorDataService,
	userContextService UserContextDataService,
	responsesStore RagResponsesRepository,
) (*SectorRag, error) {
	rag := SectorRag{
		dataService:        sectorDataService,
		userContextService: userContextService,
	}
	rag.llm = llm
	rag.topic = SECTORS
	rag.responseStore = responsesStore

	return &rag, nil
}

func (rag SectorRag) createRagContext(sectorName string) (string, error) {
	var ragContext string
	sectors, err := rag.dataService.GetSectors()
	if err != nil {
		return ragContext, &DataServiceError{Message: fmt.Sprintf("GetSectors failed: %s", err)}
	}

	for i := 0; i < len(sectors); i++ {
		sector := sectors[i]
		sectorStocks, err := rag.dataService.GetSectorStocks(sector.UrlName)
		if err != nil {
			return ragContext, &DataServiceError{Message: fmt.Sprintf("GetSectorStocks failed: %s", err)}
		}

		if sectorName == "" {
			// Keep only the top 5 stocks in case we have all sectors in the prompt
			context := sectorContext{
				sector:       sector,
				sectorStocks: sectorStocks[:5],
			}
			ragContext += fmt.Sprintf("%+v\n", context)
		} else if sectorName == sector.Name {
			context := sectorContext{
				sector:       sector,
				sectorStocks: sectorStocks,
			}
			ragContext += fmt.Sprintf("%+v\n", context)
			return ragContext, nil
		}
	}

	return ragContext, nil
}

func (rag SectorRag) GenerateRagResponse(conversation []Message, tags Tags, responseChannel chan<- string) error {
	// Format the prompt to contain the neccessary context
	ragContext, err := rag.createRagContext(tags.SectorName)
	if err != nil {
		return err
	}

	var userContext domain.UserContext
	if tags.UserID != "" {
		userContext, err = rag.userContextService.GetUserContext(tags.UserID)
		if err != nil {
			return err
		}
	}

	prompt := fmt.Sprintf(prompts.SectorsPrompt, ragContext, userContext)

	return rag.GenerateLllmResponse(prompt, conversation, responseChannel)
}
