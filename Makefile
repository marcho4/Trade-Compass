.PHONY: build financial-data parser auth frontend ai restart

build:
	docker compose build --no-cache

financial-data:
	docker compose build financial-data && docker down financial-data && docker compose up -d

parser:
	docker compose build parser && docker compose down parser && docker compose up -d

auth:
	docker compose build auth-service && docker down auth-service && docker compose up -d

frontend:
	docker compose build frontend && docker compose down frontend && docker compose up -d

ai:
	docker compose build ai-service && docker down ai-service && docker compose up -d

restart:
	docker compose down && docker compose up -d