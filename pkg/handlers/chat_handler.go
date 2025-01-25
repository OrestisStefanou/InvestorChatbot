package handlers

type Service interface {
	GenerateResponse(topic string, sessionId string, question string, responseChannel chan<- string) error
}

type ChatHandler struct {
	chatService Service
}

func (h *ChatHandler) ServeRequest() {

}
