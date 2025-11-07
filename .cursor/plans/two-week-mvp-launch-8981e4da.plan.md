<!-- 8981e4da-c70f-481b-ab6b-d23ee728bfaf 5ab88d38-3259-439b-95bb-0e43f762864e -->
# План на первые две недели: MVP фундамент

## Контекст и цель

Цель: Запустить работающий минимальный прототип с API для получения данных по 3-5 компаниям (Сбер, Яндекс, Лукойл, Газпром, ВТБ). Данные вводятся вручную, парсер пока в стороне. Фокус на качественной архитектуре и быстром старте.

**Важно**: Работаешь соло, поэтому приоритет на быстрых победах и избегании переусложнения.

---

## Неделя 1: Инфраструктура + База данных + Go сервис (скелет)

### День 1-2: Настройка окружения и структуры проекта

**Задачи:**

- Создать структуру монорепозитория:
  - `/backend` (Go сервис - metrics-service)
  - `/frontend` (Next.js - пока не трогаем)
  - `/parser` (уже есть)
  - `/migrations` (SQL миграции)
  - `/docs` (уже есть)

- Настроить Docker Compose для локальной разработки:
  - PostgreSQL 16
  - Redis (опционально, для будущего кэширования)
  - Adminer или pgAdmin для визуализации БД

- Создать `.env.example` с переменными:
  ```
  POSTGRES_HOST=localhost
  POSTGRES_PORT=5432
  POSTGRES_DB=bullrun
  POSTGRES_USER=...
  POSTGRES_PASSWORD=...
  MOEX_API_URL=https://iss.moex.com
  ```

- Создать Makefile с командами:
  - `make dev` - запуск всех сервисов
  - `make migrate-up` - применить миграции
  - `make migrate-down` - откатить миграции
  - `make test` - запуск тестов
  - `make lint` - линтеры

**Результат**: Полностью готовое окружение для разработки, можно поднять БД и начать писать код.

---

### День 3-4: Схема базы данных и миграции

**Задачи:**

- Доработать схему из `docs/database.md`:
  - Добавить недостающие поля в `companies`: `name`, `description`, `website`, `logo_url`, `market_cap`
  - Создать таблицу `industries` (отрасли внутри секторов)
  - Добавить `created_at`, `updated_at` для всех таблиц
  - Добавить индексы на часто используемые поля

- Создать SQL миграции через `golang-migrate`:
  - `001_create_sectors.up.sql`
  - `002_create_industries.up.sql`
  - `003_create_companies.up.sql`
  - `004_create_reports.up.sql`
  - `005_create_metrics.up.sql`
  - `006_create_indicators.up.sql`
  - `007_create_dividends.up.sql`
  - `008_create_users_auth.up.sql` (базовая структура)
  - `009_create_ai_analyses.up.sql`
  - `010_create_portfolios.up.sql`

- Применить миграции локально

- Заполнить справочники:
  - `sectors`: Финансы, Энергетика, IT, Потребительский сектор, и т.д.
  - `industries`: Банки, Нефть и газ, Интернет-сервисы, и т.д.

**Результат**: Рабочая схема БД, все таблицы созданы, справочники заполнены.

---

### День 5-7: Go Backend - Базовая структура и первые эндпоинты

**Задачи:**

**Структура проекта** (Clean Architecture):

```
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/          # Модели (structs)
│   ├── repository/      # Работа с БД
│   ├── usecase/         # Бизнес-логика
│   ├── handler/         # HTTP handlers
│   └── middleware/      # Middleware
├── pkg/
│   ├── database/        # Подключение к БД
│   └── logger/          # Логирование
└── go.mod
```

**Технологический стек:**

- Фреймворк: `Gin` (простой и быстрый)
- БД драйвер: `pgx/v5` (высокая производительность)
- Миграции: `golang-migrate`
- Конфиг: `viper` или простые env переменные
- Логирование: `slog` (стандартная библиотека Go 1.21+)

**Реализовать:**

1. Подключение к PostgreSQL через pgx
2. Создать модели (domain):

   - `Company`
   - `Sector`
   - `Industry`
   - `Report`
   - `Indicator`
   - `Metric`

3. Repository слой:

   - `CompanyRepository` с методами:
     - `GetAll(ctx context.Context) ([]Company, error)`
     - `GetByTicker(ctx context.Context, ticker string) (*Company, error)`
     - `GetBySector(ctx context.Context, sectorID int) ([]Company, error)`

4. HTTP handlers:

   - `GET /api/v1/health` - проверка здоровья сервиса
   - `GET /api/v1/companies` - список всех компаний
   - `GET /api/v1/companies/:ticker` - детали компании
   - `GET /api/v1/sectors` - список секторов
   - `GET /api/v1/industries` - список отраслей

5. Middleware:

   - CORS
   - Request logging
   - Recovery (panic handling)

**Результат**: Запущенный Go сервис на порту 8080, отвечающий на базовые запросы. Пока без реальных данных по компаниям (вернет пустые списки).

---

## Неделя 2: Наполнение данными + API для метрик

### День 8-9: Ручной ввод данных для 3-5 компаний

**Задачи:**

Для каждой компании (Сбер, Яндекс, Лукойл, Газпром, ВТБ):

1. Найти последние финансовые отчеты:

   - Годовой отчет (2023 или 2024)
   - Последние квартальные отчеты (Q1-Q4 2024)

2. Создать SQL скрипт для вставки данных:

   - Компания: `INSERT INTO companies (ticker, inn, name, sector_id, industry_id, ...)`
   - Отчеты: `INSERT INTO reports (company_id, report_year, report_period, ...)`
   - Метрики: `INSERT INTO metrics (report_id, revenue, net_profit, total_assets, ...)`

3. Вручную извлечь данные из отчетов:

   - Из баланса: активы, обязательства, капитал, денежные средства
   - Из P&L: выручка, валовая прибыль, операционная прибыль, чистая прибыль
   - Из Cash Flow: операционный CF, CAPEX

4. Применить скрипты вставки данных

5. Добавить тестовые данные по дивидендам (если применимо)

**Результат**: В БД есть 3-5 компаний с финансовыми данными за последние 4-8 кварталов.

---

### День 10-11: Сервис расчета индикаторов

**Задачи:**

1. Создать калькулятор метрик в `internal/calculator/`:

   - `CalculatePE(price, eps float64) float64`
   - `CalculatePB(price, bookValue float64) float64`
   - `CalculateROE(netProfit, equity float64) float64`
   - `CalculateROA(netProfit, assets float64) float64`
   - `CalculateDebtToEquity(debt, equity float64) float64`
   - `CalculateCurrentRatio(currentAssets, currentLiabilities float64) float64`
   - ... (все метрики из `docs/business.md`)

2. Создать usecase для расчета индикаторов:

   - `CalculateIndicatorsForReport(ctx context.Context, reportID int) error`
   - Берет данные из `metrics`, считает все индикаторы, сохраняет в `indicators`

3. API эндпоинты:

   - `GET /api/v1/companies/:ticker/metrics` - метрики по всем периодам
   - `GET /api/v1/companies/:ticker/indicators` - рассчитанные индикаторы
   - `POST /api/v1/internal/calculate-indicators/:reportID` - пересчет (для админа)

4. Запустить расчет для всех введенных отчетов

**Результат**: Все индикаторы автоматически рассчитываются и доступны через API.

---

### День 12-13: Интеграция с MOEX API + История цен

**Задачи:**

1. Изучить MOEX ISS API:

   - Документация: https://iss.moex.com/iss/reference/
   - Endpoint для текущей цены: `/iss/engines/stock/markets/shares/securities/{ticker}.json`
   - Endpoint для истории: `/iss/history/engines/stock/markets/shares/securities/{ticker}.json`

2. Создать MOEX клиент в `pkg/moex/`:

   - `GetCurrentPrice(ctx context.Context, ticker string) (*StockPrice, error)`
   - `GetHistoricalPrices(ctx context.Context, ticker string, from, to time.Time) ([]StockPrice, error)`
   - Rate limiting (не более 1 запроса в секунду)

3. Создать таблицу `stock_prices`:
   ```sql
   CREATE TABLE stock_prices (
     id SERIAL PRIMARY KEY,
     company_id INT REFERENCES companies(id),
     price DECIMAL(10, 2),
     volume BIGINT,
     date DATE,
     created_at TIMESTAMPTZ DEFAULT NOW()
   );
   ```

4. API эндпоинты:

   - `GET /api/v1/companies/:ticker/price` - текущая цена
   - `GET /api/v1/companies/:ticker/price-history?from=&to=` - история

5. Создать cronjob (пока просто скрипт, запускаемый вручную):

   - Обновляет цены для всех компаний
   - Сохраняет в `stock_prices`

**Результат**: Можем получать текущие цены и историю через API, данные сохраняются в БД.

---

### День 14: Тестирование, документация, итоги

**Задачи:**

1. Написать базовые unit-тесты:

   - Тесты для калькулятора метрик
   - Тесты для repository методов
   - Минимум 50% coverage критичной логики

2. Создать Swagger документацию (опционально):

   - Использовать `swaggo/swag`
   - Документировать все эндпоинты

3. Обновить документацию в `/docs`:

   - Добавить инструкцию по запуску проекта
   - Описать API эндпоинты
   - Обновить `roadmap.md` с актуальным прогрессом

4. Проверить работу всей системы end-to-end:

   - Запустить `make dev`
   - Проверить все эндпоинты через Postman или curl
   - Убедиться, что метрики считаются корректно
   - Проверить получение цен с MOEX

5. Зафиксировать достижения:

   - Создать git tag `v0.1.0-mvp-backend`
   - Подготовить краткий отчет о проделанной работе

**Результат**: Работающий MVP бэкенда с API для получения данных по компаниям, рассчитанными метриками и интеграцией с MOEX.

---

## Что НЕ делаем в эти две недели

❌ Фронтенд (Next.js) - это следующая итерация

❌ AI интеграция (Claude) - Phase 3

❌ Авторизация и монетизация - Phase 4

❌ Автоматический парсер отчетов - Phase 5

❌ Продакшн деплой - позже

---

## Критерии успеха

✅ Поднимается локально за одну команду `make dev`

✅ API возвращает данные по 3-5 компаниям

✅ Все индикаторы считаются автоматически

✅ Интеграция с MOEX работает

✅ Код покрыт базовыми тестами

✅ Документация актуальна

---

## Риски и митигация

**Риск 1**: Ручной ввод данных займет больше времени

- **Митигация**: Начать с 3 компаний, остальные добавить позже

**Риск 2**: Сложности с MOEX API (лимиты, формат данных)

- **Митигация**: Начать с простых запросов, добавить retry и кэширование

**Риск 3**: Неправильные формулы расчета метрик

- **Митигация**: Сверить с SimplyWallSt и другими источниками, написать тесты

---

## После двух недель

**Следующие шаги (Неделя 3-4):**

1. Запуск фронтенда (Next.js + Shadcn)
2. Визуализация данных по компаниям
3. Первые графики и таблицы
4. Адаптивный дизайн

**Цель месяца**: Работающий MVP, который можно показать первым тестовым пользователям.

### To-dos

- [ ] Настроить монорепозиторий, Docker Compose, Makefile и .env конфигурацию
- [ ] Создать SQL миграции для всех таблиц и заполнить справочники (sectors, industries)
- [ ] Создать структуру Go сервиса (Clean Architecture) с базовыми эндпоинтами
- [ ] Вручную ввести финансовые данные для 3-5 компаний (Сбер, Яндекс, Лукойл)
- [ ] Реализовать калькулятор всех индикаторов (P/E, ROE, ROA и т.д.) и API для метрик
- [ ] Интегрировать MOEX API для получения текущих и исторических цен акций
- [ ] Написать unit-тесты, обновить документацию и провести E2E проверку системы