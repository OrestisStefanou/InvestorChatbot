tests:
	go test -v ./...

run_investbot:
	go run cmd/investbot/main.go

build_investbot:
	go build cmd/investbot/main.go

install:
	go mod tidy
	go mod download
