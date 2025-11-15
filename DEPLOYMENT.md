# Развертывание Bull Run

Инструкция по развертыванию приложения Bull Run на сервере.

## Предварительные требования

- Docker Engine 24.0+
- Docker Compose V2
- Сервер с минимум 4GB RAM
- Открытый порт 80

## Быстрый старт

### 1. Клонирование репозитория

```bash
git clone <repository-url>
cd bull-run
```

### 2. Настройка переменных окружения

```bash
# Создайте .env файл из примера
cp .env.example .env

# Отредактируйте переменные окружения
nano .env
```

Обязательно измените следующие переменные:
- `POSTGRES_PASSWORD` - пароль для PostgreSQL
- `JWT_SECRET` - секретный ключ для JWT токенов

### 3. Запуск в Production режиме

```bash
# Сборка и запуск всех сервисов
docker compose up -d

# Проверка статуса
docker compose ps

# Просмотр логов
docker compose logs -f
```

### 4. Проверка работоспособности

```bash
# Health check
curl http://89.169.137.57/health

# Должен вернуть: healthy
```

Откройте в браузере: http://89.169.137.57

## Структура сервисов

- **nginx** - Reverse proxy (порт 80)
- **frontend** - Next.js приложение (внутренний порт 3000)
- **backend** - API сервер (внутренний порт 8080)
- **postgres** - База данных PostgreSQL
- **redis** - Кеш и хранилище сессий
- **parser** - Сервис парсинга данных по акциям
- **cron-jobs** - Периодические задачи

## Режимы работы

### Production (по умолчанию)

```bash
# В .env:
NODE_ENV=production
FRONTEND_DOCKERFILE=Dockerfile.prod

docker compose up -d
```

### Development

```bash
# В .env:
NODE_ENV=development
FRONTEND_DOCKERFILE=Dockerfile

docker compose up -d
```

## Управление сервисами

```bash
# Остановка всех сервисов
docker compose down

# Остановка с удалением volumes
docker compose down -v

# Пересборка и перезапуск
docker compose up -d --build

# Рестарт конкретного сервиса
docker compose restart nginx

# Просмотр логов конкретного сервиса
docker compose logs -f frontend
```

## Обновление приложения

```bash
# Остановка сервисов
docker compose down

# Получение обновлений
git pull

# Пересборка и запуск
docker compose up -d --build
```

## Резервное копирование

### База данных

```bash
# Создание бэкапа
docker compose exec postgres pg_dump -U postgres bullrun > backup_$(date +%Y%m%d_%H%M%S).sql

# Восстановление
docker compose exec -T postgres psql -U postgres bullrun < backup.sql
```

## Мониторинг

### Использование ресурсов

```bash
# Статистика контейнеров
docker stats

# Размер volumes
docker system df -v
```

### Логи

```bash
# Все логи
docker compose logs -f

# Последние 100 строк
docker compose logs --tail=100

# Логи конкретного сервиса
docker compose logs -f nginx
docker compose logs -f frontend
docker compose logs -f backend
```

## Troubleshooting

### Контейнер не запускается

```bash
# Проверка логов
docker compose logs <service-name>

# Проверка конфигурации
docker compose config
```

### Проблемы с сетью

```bash
# Пересоздание сети
docker compose down
docker network prune
docker compose up -d
```

### Очистка системы

```bash
# Удаление неиспользуемых образов и контейнеров
docker system prune -a

# Удаление volumes (ОСТОРОЖНО! Потеря данных)
docker volume prune
```

## Безопасность

1. **Измените пароли** в .env файле перед production развертыванием
2. **Настройте firewall** для ограничения доступа к портам
3. **Используйте HTTPS** в production (настройте SSL/TLS)
4. **Регулярно обновляйте** Docker образы
5. **Настройте мониторинг** и алерты

## SSL/TLS (HTTPS)

Для настройки HTTPS:

1. Получите SSL сертификат (Let's Encrypt)
2. Обновите nginx конфигурацию
3. Добавьте volumes для сертификатов в docker-compose.yml

Пример будет добавлен позже.

## Поддержка

При возникновении проблем:
1. Проверьте логи: `docker compose logs -f`
2. Проверьте статус: `docker compose ps`
3. Создайте issue в репозитории
