package services

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

type Llm interface {
	GenerateResponse(conversation []Message, responseChannel chan<- string) error
}
