include .env.local

migrate-up:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL} up

migrate-down:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL} down

generate:
	@sqlc generate
	
