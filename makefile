GOOSE := goose
DB_URL := postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
MIGRATIONS_DIR := sql/schema

.PHONY: up down status sqlc

## Run all migrations
up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

## Roll back the last migration
down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

## Show migration status
status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

## Generate sqlc code
sqlc:
	sqlc generate
