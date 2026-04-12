.PHONY: build financial-data parser auth frontend ai restart

build:
	docker compose build --no-cache

financial-data:
	docker compose up -d --build --force-recreate financial-data

parser:
	docker compose up -d --build --force-recreate parser

auth:
	docker compose up -d --build --force-recreate auth-service

frontend:
	docker compose build frontend
	docker compose up -d --force-recreate frontend

ai:
	docker compose up -d --build --force-recreate ai-service

restart:
	docker compose down && docker compose up -d
