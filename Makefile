.PHONY: build
build:
	@go build -o bin/api ./cmd/api/main.go
	@go build -o bin/migrate ./cmd/migrate/main.go

.PHONY: dev
dev:
	@go run ./cmd/api/main.go

.PHONY: migrate
migrate:
	@go run ./cmd/migrate/main.go

all: run
