.PHONY: up down build run migrate-up migrate-down test lint

up:
	docker compose up --build

down:
	docker compose down

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

migrate-up:
	docker compose run --rm migrate

migrate-down:
	@set -a; . ./.env 2>/dev/null; set +a; \
	docker compose run --rm migrate \
		-path /migrations \
		-database "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@postgres:5432/$${POSTGRES_DB}?sslmode=disable" \
		down 1

test:
	go test ./...

lint:
	go vet ./...
