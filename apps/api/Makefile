# This is the alternate way to migrate and generate SQLC queries
#
# Note: Please have an .env.development in the root of the api/ folder

include .env.local

migrate-up:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL} up

migrate-down:
	@goose -dir ./internal/db/migrations postgres ${DATABASE_URL} down

generate:
	@sqlc generate
	
