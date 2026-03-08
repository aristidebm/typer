.PHONY: format test build

run:
	@go run cmd/main.go

format:
	@go fmt ./...

test:
	@go test -v 

build:
	@go build -o build/typer cmd/main.go
