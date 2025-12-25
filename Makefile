.PHONY: api
api:
	@go build -o bin/api ./cmd/api/main.go

.PHONY: dev
dev:
	@go run ./cmd/api/main.go

all: run
