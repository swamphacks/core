include .env.dev

migrate-up:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL_MIGRATION} up

migrate-down:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL_MIGRATION} down

generate:
	@sqlc generate
	
generate-openapi:
	swag init --dir cmd/api,internal/api/handlers --parseDependency -v3.1