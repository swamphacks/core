include .env.dev

migrate-up:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL_MIGRATION} up

migrate-down:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL_MIGRATION} down

# Use this command when developing inside a dev container.
# Maybe we could figure out a way to connect directly to the host's localhost, instead of using docker network to connect to the postgres db
migrate-up-devcontainer:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL} up

migrate-down-devcontainer:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL} down

generate:
	@sqlc generate
	
openapi-generate:
	swag init --dir cmd/api,internal/api/handlers --parseDependency --requiredByDefault -v3.1

openapi-format:
	swag f --dir cmd/api,internal/api/handlers
