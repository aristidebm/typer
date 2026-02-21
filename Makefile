.PHONY: format test build

format:
	@go fmt ./...

test:
	@go test -v 

run:
	@go run cmd/main.go

build:
	@go build -o build/typer cmd/main.go
