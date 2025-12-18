.PHONY: deps run db-up db-down migrate-up migrate-down

deps:
	go mod tidy

db-up:
	docker-compose up -d postgres

db-down:
	docker-compose down -v

run:
	go run ./cmd/api

migrate-up:
	migrate -path ./migrations -database "$$DATABASE_URL" up

migrate-down:
	migrate -path ./migrations -database "$$DATABASE_URL" down 1
