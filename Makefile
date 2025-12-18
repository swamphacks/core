local:
	docker compose up

api:
	docker compose up api

bat:
	docker compose up api bat_worker asynqmon

storage:
	docker compose up postgres redis

backend:
	docker compose up api email_worker bat_worker asynqmon
