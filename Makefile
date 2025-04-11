tests:
	go test -v ./...

run_investbot:
	go run cmd/investbot/main.go

install:
	go mod tidy
	go mod download
