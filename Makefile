local:
	docker compose up

api:
	docker compose up api

storage:
	docker compose up postgres redis
