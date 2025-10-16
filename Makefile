local:
	docker compose up

api:
	docker compose up api

storage:
	docker compose up postgres redis

backend:
	docker compose up api email_worker asynqmon
