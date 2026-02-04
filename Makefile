build:
	docker compose build --no-cache

rebuild-financial-data:
	docker compose build financial-data --no-cache && docker down financial-data && docker compose up -d

restart:
	docker compose down && docker compose up -d