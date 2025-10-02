include .env.dev

migrate-up:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL_MIGRATION} up

migrate-down:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL_MIGRATION} down

generate:
	@sqlc generate
	
openapi-generate:
	swag init --dir cmd/api,internal/api/handlers --parseDependency --requiredByDefault -v3.1

openapi-format:
	swag f --dir cmd/api,internal/api/handlers
