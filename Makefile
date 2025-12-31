.PHONY: build
build:
	@go build -o bin/api ./cmd/api/main.go
	@go build -o bin/migrate ./cmd/migrate/main.go
	@go build -o bin/test ./cmd/test/main.go

.PHONY: dev
dev:
	@air

.PHONY: migrate
migrate:
	@go run ./cmd/migrate/main.go

.PHONY: test
test: build
	@./bin/test

all: run
