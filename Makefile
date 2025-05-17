generate-openapi:
	swag init -g apps/api/cmd/api/main.go -o shared/openapi
	cp shared/openapi/docs.go apps/api/internal/api/docs
