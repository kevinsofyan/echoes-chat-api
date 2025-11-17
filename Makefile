include .env
export $(shell sed 's/=.*//' .env)

run:
	go run cmd/server/main.go

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

tidy:
	go mod tidy

test:
	go test -v ./...

test-coverage:
	go test -cover ./...

test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

migrate-up:
	migrate -path migrations -database "${DATABASE_URL}" up

migrate-down:
	migrate -path migrations -database "${DATABASE_URL}" down

migrate-up-local:
	migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}&search_path=${DB_SCHEMA}" up

migrate-down-local:
	migrate -path migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}&search_path=${DB_SCHEMA}" down

migrate-create:
	@read -p "Enter migration name: " name; \
    migrate create -ext sql -dir migrations -seq $$name

db-seed:
	docker-compose exec db psql -U ${DB_USER} -d ${DB_NAME} -c "SET search_path TO ${DB_SCHEMA};" -f /seeds/seed.sql

swagger:
	/Users/kevinsofyan/go/bin/swag init -g cmd/server/main.go

.PHONY: run docker-up docker-down tidy test test-coverage test-coverage-html migrate-up migrate-down migrate-up-local migrate-down-local migrate-create db-seed swagger