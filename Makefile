build:
	docker compose build --no-cache

financial-data:
	docker compose build financial-data --no-cache && docker down financial-data && docker compose up -d

parser:
	docker compose build parser --no-cache && docker down parser && docker compose up -d

auth:
	docker compose build auth-service --no-cache && docker down auth-service && docker compose up -d

frontend:
	docker compose build frontend --no-cache && docker down frontend && docker compose up -d

ai:
	docker compose build ai-service --no-cache && docker down ai-service && docker compose up -d

restart:
	docker compose down && docker compose up -d