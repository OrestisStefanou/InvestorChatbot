package services

type Llm interface {
	GenerateResponse(conversation []map[string]string, responseChannel chan<- string) error
}
