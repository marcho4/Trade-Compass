# Backend Service

Placeholder для backend сервиса.

## TODO

Реализовать следующий функционал:

### API Endpoints

- Авторизация и регистрация
  - POST /api/auth/register
  - POST /api/auth/login
  - POST /api/auth/refresh
  - POST /api/auth/logout
  - POST /api/auth/verify-email
  - POST /api/auth/forgot-password
  - POST /api/auth/reset-password

- OAuth
  - GET /api/auth/oauth/google
  - GET /api/auth/oauth/google/callback
  - GET /api/auth/oauth/yandex
  - GET /api/auth/oauth/yandex/callback

- Пользователь
  - GET /api/users/me
  - PUT /api/users/me
  - DELETE /api/users/me

- Акции и аналитика
  - GET /api/stocks
  - GET /api/stocks/:ticker
  - GET /api/stocks/:ticker/analysis

- AI Assistant
  - POST /api/ai/chat

## Технологический стек

Рекомендуется:
- Go + Gin/Echo или Node.js + Express/Fastify
- PostgreSQL для хранения данных
- Redis для кеширования и сессий
- JWT для авторизации
