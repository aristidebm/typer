.PHONY: format test

format:
	@go fmt ./...

test:
	@go test -v 

run:
	@go run cmd/main.go
