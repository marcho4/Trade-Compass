.PHONY: build financial-data parser auth frontend ai restart

build:
	docker compose build --no-cache

financial-data:
	docker compose up -d --build --force-recreate financial-data
	docker compose restart nginx

parser:
	docker compose up -d --build --force-recreate parser
	docker compose restart nginx

auth:
	docker compose up -d --build --force-recreate auth-service
	docker compose restart nginx

frontend:
	docker compose up -d --build --force-recreate frontend
	docker compose restart nginx

ai:
	docker compose up -d --build --force-recreate ai-service
	docker compose restart nginx

restart:
	docker compose down && docker compose up -d
