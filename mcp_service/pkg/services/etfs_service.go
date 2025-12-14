package services

import (
	"fmt"
	"investbot/pkg/domain"
	"investbot/pkg/services/prompts"
	"strings"
)

type EtfDataService interface {
	GetEtfs() ([]domain.Etf, error)
	GetEtfOverview(symbol string) (domain.EtfOverview, error)
}

type etfContext struct {
	etf domain.Etf
}

type etfOverviewContext struct {
	etfOverview domain.EtfOverview
}

type EtfRag struct {
	BaseRag
	dataService        EtfDataService
	userContextService UserContextDataService
}

func NewEtfRag(
	llm Llm,
	etfDataService EtfDataService,
	userContextService UserContextDataService,
	responsesStore RagResponsesRepository,
) (*EtfRag, error) {
	rag := EtfRag{
		dataService:        etfDataService,
		userContextService: userContextService,
	}
	rag.llm = llm
	rag.topic = ETFS
	rag.responseStore = responsesStore

	return &rag, nil
}

func (rag EtfRag) createRagContext(etfSymbols []string) (string, error) {
	var ragContext string

	if len(etfSymbols) > 0 {
		for _, etfSymbol := range etfSymbols {
			etfOverview, err := rag.dataService.GetEtfOverview(etfSymbol)
			if err != nil {
				return ragContext, &DataServiceError{Message: fmt.Sprintf("GetEtfOverview failed: %s", err)}
			}
			context := etfOverviewContext{etfOverview: etfOverview}
			ragContext += fmt.Sprintf("%+v\n", context)
		}
		return ragContext, nil
	}

	etfs, err := rag.dataService.GetEtfs()
	if err != nil {
		return ragContext, &DataServiceError{Message: fmt.Sprintf("GetEtfs failed: %s", err)}
	}

	var context etfContext
	for _, etf := range etfs {
		// Because of the huge number of ETFs we keep only the ones
		// with aum higher than 2.5B
		if etf.Aum > 2500000000 {
			context.etf = etf
			ragContext += fmt.Sprintf("%+v\n", context)
		}
	}

	return ragContext, nil
}

func (rag EtfRag) GenerateRagResponse(conversation []Message, tags Tags, responseChannel chan<- string) error {
	// Format the prompt to contain the neccessary context
	ragContext, err := rag.createRagContext(tags.EtfSymbols)
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

	prompt := fmt.Sprintf(prompts.EtfsPrompt, ragContext, userContext)

	return rag.GenerateLllmResponse(prompt, conversation, responseChannel)
}

type EtfService struct {
	dataService EtfDataService
}

func NewEtfService(dataService EtfDataService) (*EtfService, error) {
	return &EtfService{dataService: dataService}, nil
}

type EtfFilterOptions struct {
	SearchString string
}

func (f EtfFilterOptions) IsEmpty() bool {
	return f.SearchString == ""
}

func (f EtfFilterOptions) HasSearchString() bool {
	return f.SearchString != ""
}

func (s EtfService) GetEtfs(filters EtfFilterOptions) ([]domain.Etf, error) {
	etfs, err := s.dataService.GetEtfs()
	if err != nil {
		return nil, err
	}

	if filters.IsEmpty() {
		return etfs, nil
	}

	if filters.HasSearchString() {
		filteredEtfs := make([]domain.Etf, 0)
		for _, e := range etfs {
			search := strings.ToLower(filters.SearchString)
			symbol := strings.ToLower(e.Symbol)
			etfName := strings.ToLower(e.Name)
			if strings.Contains(symbol, search) || strings.Contains(etfName, search) {
				filteredEtfs = append(filteredEtfs, e)
			}

		}
		return filteredEtfs, nil
	}

	return etfs, nil
}

func (s EtfService) GetEtf(etfSymbol string) (domain.EtfOverview, error) {
	return s.dataService.GetEtfOverview(etfSymbol)
}
