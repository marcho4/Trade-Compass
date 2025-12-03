# Авторизация и аутентификация

## Обзор

Система авторизации построена на JWT токенах, хранящихся в HttpOnly cookies. Это обеспечивает безопасность от XSS атак и автоматическую отправку токенов при каждом запросе.

## Архитектура

```
┌─────────────────────────────────────────────────────────────────────┐
│                           BROWSER                                    │
│  ┌──────────────┐     ┌───────────────┐     ┌───────────────┐      │
│  │ React App    │     │ HttpOnly      │     │ middleware.ts │      │
│  │              │────▶│ Cookies       │◀────│ (проверка)    │      │
│  │ useAuth()    │     │ • accessToken │     └───────────────┘      │
│  │ api-client   │     │ • refreshToken│                             │
│  └──────┬───────┘     └───────┬───────┘                             │
└─────────┼─────────────────────┼─────────────────────────────────────┘
          │ credentials:'include'│
          ▼                     ▼ (автоматически)
┌─────────────────────────────────────────────────────────────────────┐
│                           NGINX                                      │
│                    /api/auth/* → auth-service                        │
└─────────────────────────────────────────────────────────────────────┘
          │
          ▼
┌─────────────────────────────────────────────────────────────────────┐
│                       AUTH-SERVICE (Go)                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────────┐ │
│  │ handlers.go │  │ service.go  │  │ jwt_service.go              │ │
│  │ SetCookie() │  │ RefreshTok  │  │ • GenerateAccessToken       │ │
│  │ HttpOnly    │  │ Rotation    │  │ • 15min / 15days TTL        │ │
│  └─────────────┘  └─────────────┘  └─────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
          │
          ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        PostgreSQL                                    │
│   refresh_tokens: id, token_hash, user_id, device_info, expires_at  │
└─────────────────────────────────────────────────────────────────────┘
```

## JWT Токены

### Access Token
- **TTL:** 15 минут
- **Хранение:** HttpOnly cookie `accessToken`
- **Содержит:** userId, name, status

### Refresh Token
- **TTL:** 15 дней
- **Хранение:** HttpOnly cookie `refreshToken`
- **Содержит:** userId, tokenId (для отзыва)
- **Ротация:** При каждом refresh старый токен аннулируется

## Cookies

Все cookies устанавливаются с флагами:

```go
http.SetCookie(w, &http.Cookie{
    Name:     "accessToken",
    Value:    accessToken,
    Path:     "/",
    Domain:   domain,
    MaxAge:   accessTokenMaxAge,
    HttpOnly: true,              // Защита от XSS
    Secure:   isSecure,          // Только HTTPS в production
    SameSite: http.SameSiteLaxMode, // Защита от CSRF
})
```

## API Endpoints (auth-service)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/register` | Регистрация нового пользователя |
| POST | `/login` | Авторизация по email/password |
| POST | `/refresh` | Обновление токенов |
| POST | `/logout` | Выход из системы |
| GET | `/me` | Получение данных текущего пользователя |
| GET | `/yandex/login` | Получение URL для OAuth через Яндекс |
| GET | `/callback/yandex` | Callback от Яндекс OAuth |

## OAuth через Яндекс

### Flow

1. Пользователь нажимает "Войти через Яндекс"
2. Frontend запрашивает `/auth/yandex/login`
3. Backend генерирует `state`, сохраняет в cookie, возвращает URL
4. Frontend редиректит на Яндекс
5. Пользователь авторизуется в Яндексе
6. Яндекс редиректит на `/callback/yandex?code=...&state=...`
7. Backend проверяет `state` (CSRF защита)
8. Backend обменивает `code` на access token
9. Backend получает данные пользователя от Яндекса
10. Backend создаёт/находит пользователя, генерирует JWT
11. Backend устанавливает cookies, редиректит на `/welcome`

### State параметр (CSRF защита)

```go
// Генерация
state, _ := generateSecureState() // 32 рандомных байта в base64

// Сохранение в cookie
http.SetCookie(w, &http.Cookie{
    Name:     "oauth_state",
    Value:    state,
    MaxAge:   300, // 5 минут
    HttpOnly: true,
})

// Проверка при callback
savedState, _ := r.Cookie("oauth_state")
if savedState.Value != r.URL.Query().Get("state") {
    // CSRF атака!
}
```

## Frontend

### Файлы

| Файл | Описание |
|------|----------|
| `lib/api-client.ts` | API клиент с авто-refresh токенов |
| `contexts/AuthContext.tsx` | React Context для состояния авторизации |
| `middleware.ts` | Next.js middleware для защиты роутов |
| `components/auth/YandexAuthButton.tsx` | Кнопка OAuth через Яндекс |
| `components/providers/Providers.tsx` | Обёртка с AuthProvider |

### Использование useAuth()

```tsx
import { useAuth } from "@/contexts/AuthContext";

function MyComponent() {
  const { 
    user,           // Данные пользователя или null
    isLoading,      // Идёт загрузка
    isAuthenticated,// Авторизован ли
    login,          // Функция входа
    register,       // Функция регистрации
    logout,         // Функция выхода
    refreshUser     // Обновить данные пользователя
  } = useAuth();

  // Пример входа
  const handleLogin = async () => {
    try {
      await login(email, password);
      router.push('/dashboard');
    } catch (error) {
      setError(error.message);
    }
  };
}
```

### Авто-refresh токенов

API клиент автоматически обновляет токены при 401 ошибке:

```typescript
// При получении 401:
// 1. Отправляет POST /auth/refresh
// 2. Если успешно — повторяет оригинальный запрос
// 3. Если нет — редиректит на /auth
```

### Middleware (защита роутов)

```typescript
// Защищённые роуты (требуют авторизации)
const PROTECTED_ROUTES = ["/dashboard"];

// Auth роуты (редирект если уже авторизован)
const AUTH_ROUTES = ["/auth", "/auth/register"];
```

## Переменные окружения

### Backend (auth-service)

```env
JWT_SECRET=your-secret-key-at-least-32-bytes
YANDEX_CLIENT_ID=your-yandex-client-id
YANDEX_CLIENT_SECRET=your-yandex-client-secret
FRONTEND_URL=https://trade-compass.ru
COOKIE_DOMAIN=.trade-compass.ru
```

### Frontend

```env
NEXT_PUBLIC_API_URL=/api
```

## Безопасность

### Защита от XSS
- Токены в HttpOnly cookies (JavaScript не может прочитать)

### Защита от CSRF
- SameSite=Lax на cookies
- State параметр в OAuth flow

### Защита от кражи refresh token
- Refresh tokens хранятся в БД
- При refresh старый токен аннулируется (rotation)
- Можно отозвать все сессии пользователя

### Production чек-лист

- [ ] `Secure: true` на cookies (только HTTPS)
- [ ] `COOKIE_DOMAIN` настроен правильно
- [ ] `JWT_SECRET` — минимум 32 байта, криптографически стойкий
- [ ] Rate limiting на `/login`, `/register`, `/refresh`
- [ ] Логирование подозрительной активности

## Расширение

### Добавление нового OAuth провайдера

1. Добавить конфиг в `config.go`:
```go
type GoogleOAuthConfig struct {
    ClientID     string
    ClientSecret string
}
```

2. Добавить handler в `handlers.go`:
```go
func (h *Handlers) HandleGoogleLogin(w http.ResponseWriter, r *http.Request)
func (h *Handlers) HandleGoogleCallback(w http.ResponseWriter, r *http.Request)
```

3. Добавить роуты в `server.go`:
```go
s.Router.Get("/google/login", s.Handlers.HandleGoogleLogin)
s.Router.Get("/callback/google", s.Handlers.HandleGoogleCallback)
```

4. Добавить кнопку на фронте:
```tsx
<GoogleAuthButton />
```

