package services

type Llm interface {
	GenerateResponse(conversation []map[string]string, responseChannel chan<- string) error
}

type LLM interface {
	GenerateResponse(conversation []Message, responseChannel chan<- string) error
}

type ActorRole string

const (
	System    ActorRole = "system"
	Assistant ActorRole = "assistant"
	User      ActorRole = "user"
)

type Message struct {
	Content string
	Role    ActorRole
}
