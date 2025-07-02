include .env.dev

migrate-up:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL_MIGRATIONS} up

migrate-down:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL_MIGRATIONS} down

generate:
	@sqlc generate
	
