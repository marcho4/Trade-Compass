# Nginx Configuration

Конфигурация Nginx для приложения Bull Run.

## Структура

- `nginx.conf` - Основная конфигурация Nginx
- `default.conf` - Конфигурация виртуального хоста

## Функционал

### Reverse Proxy

Nginx работает как обратный прокси и направляет запросы:

- `/api/*` → Backend API (порт 8080)
- `/*` → Frontend Next.js (порт 3000)
- `/health` → Health check endpoint

### Безопасность

Настроены следующие заголовки безопасности:
- `X-Frame-Options: SAMEORIGIN`
- `X-Content-Type-Options: nosniff`
- `X-XSS-Protection: 1; mode=block`

### Оптимизация

- **Gzip сжатие** для текстовых файлов
- **Кеширование** статических файлов Next.js (365 дней)
- **WebSocket поддержка** для Next.js HMR в режиме разработки

## Использование

### Production режим

```bash
# Создайте .env файл из примера
cp .env.example .env

# Отредактируйте переменные окружения
nano .env

# Запустите все сервисы
docker compose up -d
```

### Development режим

```bash
# В .env установите:
# NODE_ENV=development
# FRONTEND_DOCKERFILE=Dockerfile

docker compose up -d
```

## Порты

- **80** - Nginx (HTTP)
- **3000** - Frontend (внутренний, доступ через Nginx)
- **8080** - Backend API (внутренний, доступ через Nginx)

## Логи

```bash
# Логи Nginx
docker compose logs -f nginx

# Все логи
docker compose logs -f
```

## Health Check

Nginx настроен с health check endpoint:
```bash
curl http://89.169.137.57/health
```

Ответ: `healthy`
