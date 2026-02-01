# Financial Data Service

Микросервис для управления финансовыми данными компаний.

## Возможности

- Управление финансовыми коэффициентами (ratios)
- Хранение сырых финансовых данных (raw data)
- Информация о компаниях и секторах
- Данные о дивидендах
- Макроэкономические показатели (ставки ЦБ)
- Новости компаний и секторов
- Получение котировок с MOEX

## Архитектура

Сервис построен по принципам Clean Architecture:

```
cmd/
  main.go              - Точка входа приложения
internal/
  application/         - HTTP handlers и роутинг
  domain/              - Бизнес-логика и модели
  infrastructure/      - Работа с БД, внешние API
migrations/            - Миграции базы данных
```

## Переменные окружения

### Обязательные

- `DB_URL` - Строка подключения к PostgreSQL

  ```
  postgres://user:password@localhost:5432/financial_data?sslmode=disable
  ```

### Опциональные

#### Сервер

- `SERVER_PORT` - Порт сервера (по умолчанию: `8082`)
- `SERVER_READ_TIMEOUT` - Таймаут чтения (по умолчанию: `15s`)
- `SERVER_WRITE_TIMEOUT` - Таймаут записи (по умолчанию: `15s`)
- `SERVER_IDLE_TIMEOUT` - Таймаут idle (по умолчанию: `60s`)

#### База данных

- `DB_MAX_CONNS` - Максимум соединений (по умолчанию: `25`)
- `DB_MIN_CONNS` - Минимум соединений (по умолчанию: `5`)
- `DB_MAX_CONN_LIFETIME` - Время жизни соединения (по умолчанию: `30m`)
- `DB_MAX_CONN_IDLE_TIME` - Время простоя соединения (по умолчанию: `5m`)

#### Безопасность

- `ADMIN_API_KEY` - API ключ для защищённых эндпоинтов
- `ALLOWED_ORIGINS` - Разрешённые CORS origins (по умолчанию: `*`)

## API Endpoints

### Health Check

- `GET /health` - Проверка работоспособности сервиса

### Ratios (Финансовые коэффициенты)

- `GET /ratios/{ticker}` - Получить коэффициенты по тикеру
- `GET /ratios/sector/{sector_id}` - Средние коэффициенты по сектору
- `POST /ratios/{ticker}` - Создать коэффициенты (требует API ключ)
- `PUT /ratios/{ticker}` - Обновить коэффициенты (требует API ключ)
- `DELETE /ratios/{ticker}` - Удалить коэффициенты (требует API ключ)

### Companies (Компании)

- `GET /companies` - Получить все компании
- `GET /companies/{ticker}` - Получить компанию по тикеру
- `GET /companies/sector/{sector_id}` - Получить компании по сектору
- `POST /companies` - Создать компанию (требует API ключ)
- `PUT /companies/{ticker}` - Обновить компанию (требует API ключ)
- `DELETE /companies/{ticker}` - Удалить компанию (требует API ключ)

### Raw Data (Сырые данные)

- `GET /raw-data/{ticker}/latest` - Последние данные по тикеру
- `GET /raw-data/{ticker}/history` - История данных по тикеру
- `GET /raw-data/{ticker}/{year}/{period}` - Данные за период
- `POST /raw-data` - Создать данные (требует API ключ)
- `PUT /raw-data` - Обновить данные (требует API ключ)
- `DELETE /raw-data/{ticker}/{year}/{period}` - Удалить данные (требует API ключ)

### Dividends (Дивиденды)

- `GET /dividends/{ticker}` - Получить дивиденды по тикеру
- `GET /dividends/id/{id}` - Получить дивиденд по ID
- `POST /dividends` - Создать дивиденд (требует API ключ)
- `PUT /dividends/{id}` - Обновить дивиденд (требует API ключ)
- `DELETE /dividends/{id}` - Удалить дивиденд (требует API ключ)

### Macro (Макроэкономика)

- `GET /macro/current` - Текущая ставка ЦБ
- `GET /macro/date/{date}` - Ставка на дату
- `GET /macro/history?from={from}&to={to}` - История ставок
- `POST /macro` - Создать запись (требует API ключ)
- `PUT /macro/{date}` - Обновить ставку (требует API ключ)
- `DELETE /macro/{date}` - Удалить запись (требует API ключ)

### News (Новости)

- `GET /news/{id}` - Получить новость по ID
- `GET /news/ticker/{ticker}` - Новости по тикеру
- `GET /news/sector/{sector_id}` - Новости по сектору
- `POST /news` - Создать новость (требует API ключ)
- `PUT /news/{id}` - Обновить новость (требует API ключ)
- `DELETE /news/{id}` - Удалить новость (требует API ключ)

### Price (Котировки)

- `GET /price/{ticker}?days={days}&interval={interval}` - Свечи MOEX

## Аутентификация

Для защищённых эндпоинтов (POST, PUT, DELETE) требуется заголовок:

```
X-API-Key: your-api-key-here
```

## Запуск

### Локально

```bash
# Установка зависимостей
go mod download

# Запуск миграций
export DB_URL="postgres://user:password@localhost:5432/financial_data?sslmode=disable"
export ADMIN_API_KEY="your-secret-key"

# Запуск сервиса
go run cmd/main.go
```

### Docker

```bash
docker build -t financial-data .
docker run -p 8082:8082 \
  -e DB_URL="postgres://user:password@db:5432/financial_data" \
  -e ADMIN_API_KEY="your-secret-key" \
  financial-data
```

### Docker Compose

```bash
docker-compose up
```

## Оптимизации

### База данных

- Пул соединений с настраиваемыми параметрами
- Connection pooling с min/max настройками
- Health checks для соединений
- Prepared statements где возможно

### HTTP

- Graceful shutdown
- Timeouts для чтения/записи
- Request ID для трейсинга
- Structured logging (JSON)

### Безопасность

- API key аутентификация
- CORS middleware
- Rate limiting (рекомендуется добавить на уровне nginx)
- Input validation

## Метрики и Мониторинг

Рекомендуется добавить:

- Prometheus metrics endpoint
- Distributed tracing (Jaeger/Zipkin)
- Error tracking (Sentry)

## Разработка

### Код стайл

- Следуем Effective Go
- Используем golangci-lint
- Покрытие тестами > 80%

### Миграции

```bash
# Создать новую миграцию
migrate create -ext sql -dir migrations -seq migration_name

# Применить миграции
migrate -path migrations -database $DB_URL up

# Откатить миграции
migrate -path migrations -database $DB_URL down
```

## Производительность

### Рекомендации по настройке БД

```sql
-- Индексы для частых запросов
CREATE INDEX idx_ratios_ticker ON ratios(ticker);
CREATE INDEX idx_ratios_sector ON ratios(sector);
CREATE INDEX idx_companies_sector ON companies(sector_id);
CREATE INDEX idx_news_ticker ON news(ticker);
CREATE INDEX idx_news_sector ON news(sector_id);
CREATE INDEX idx_dividends_ticker ON dividends(ticker);
```

### Настройка PostgreSQL

```ini
max_connections = 100
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 4MB
min_wal_size = 1GB
max_wal_size = 4GB
```

## Лицензия

MIT
