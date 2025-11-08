# План интеграции OAuth авторизации (Google и Яндекс)

## Обзор

Данный документ содержит детальный план интеграции OAuth 2.0 авторизации через Google и Яндекс ID для платформы Bull Run.

**Технологический стек:**
- Backend: Go
- Database: PostgreSQL
- Cache/Sessions: Redis
- Frontend: Next.js (уже готов)

**Схема БД:** Уже спроектирована, таблицы `users`, `auth`, `provider_auth`, `subscriptions` готовы.

---

## Фаза 1: Подготовка инфраструктуры

### 1.1. Создание структуры Go-проекта

**Задача:** Инициализировать Go-проект с правильной структурой папок

**Действия:**
```bash
backend/
├── cmd/
│   └── api/
│       └── main.go                    # Точка входа приложения
├── internal/
│   ├── config/
│   │   └── config.go                  # Конфигурация приложения
│   ├── database/
│   │   ├── postgres.go                # Подключение к PostgreSQL
│   │   └── redis.go                   # Подключение к Redis
│   ├── models/
│   │   ├── user.go                    # User модель
│   │   ├── auth.go                    # Auth модель
│   │   ├── provider_auth.go           # ProviderAuth модель
│   │   └── subscription.go            # Subscription модель
│   ├── repository/
│   │   ├── user_repository.go         # CRUD для users
│   │   ├── auth_repository.go         # CRUD для auth
│   │   └── provider_repository.go     # CRUD для provider_auth
│   ├── service/
│   │   ├── auth_service.go            # Бизнес-логика авторизации
│   │   ├── oauth_service.go           # OAuth логика
│   │   ├── jwt_service.go             # Работа с JWT токенами
│   │   └── user_service.go            # Бизнес-логика пользователей
│   ├── handler/
│   │   ├── auth_handler.go            # HTTP handlers для auth
│   │   └── oauth_handler.go           # HTTP handlers для OAuth
│   ├── middleware/
│   │   ├── auth_middleware.go         # JWT верификация
│   │   ├── cors_middleware.go         # CORS настройки
│   │   └── rate_limit_middleware.go   # Rate limiting
│   └── utils/
│       ├── errors.go                  # Обработка ошибок
│       ├── response.go                # Стандартизация ответов
│       └── validator.go               # Валидация данных
├── migrations/
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_auth_table.up.sql
│   ├── 002_create_auth_table.down.sql
│   ├── 003_create_provider_auth_table.up.sql
│   ├── 003_create_provider_auth_table.down.sql
│   ├── 004_create_subscriptions_table.up.sql
│   └── 004_create_subscriptions_table.down.sql
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Dockerfile
└── Makefile
```

**Зависимости (go.mod):**
```go
require (
    github.com/gin-gonic/gin v1.9.1              // HTTP фреймворк
    github.com/lib/pq v1.10.9                    // PostgreSQL драйвер
    github.com/redis/go-redis/v9 v9.5.1          // Redis клиент
    github.com/golang-jwt/jwt/v5 v5.2.0          // JWT токены
    golang.org/x/oauth2 v0.18.0                  // OAuth2 клиент
    github.com/google/uuid v1.6.0                // UUID генерация
    github.com/joho/godotenv v1.5.1              // .env файлы
    github.com/golang-migrate/migrate/v4 v4.17.0 // Миграции БД
    golang.org/x/crypto v0.21.0                  // Хэширование паролей
    github.com/go-playground/validator/v10 v10.19.0 // Валидация
)
```

**Оценка времени:** 2-3 часа

---

### 1.2. Создание миграций базы данных

**Задача:** Создать SQL миграции для таблиц авторизации

**001_create_users_table.up.sql:**
```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_users_created_at ON users(created_at);
```

**002_create_auth_table.up.sql:**
```sql
CREATE TABLE IF NOT EXISTS auth (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(100) UNIQUE NOT NULL,
    hashed_password VARCHAR(100)
);

CREATE UNIQUE INDEX idx_auth_email ON auth(email);
```

**003_create_provider_auth_table.up.sql:**
```sql
CREATE TABLE IF NOT EXISTS provider_auth (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_user_id VARCHAR(255) NOT NULL,
    provider_type VARCHAR(20) NOT NULL,
    email VARCHAR(100),
    avatar_url TEXT,
    access_token TEXT,
    refresh_token TEXT,
    token_expiry TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_provider_auth_unique ON provider_auth(provider_type, provider_user_id);
CREATE INDEX idx_provider_auth_user_id ON provider_auth(user_id);
```

**004_create_subscriptions_table.up.sql:**
```sql
CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    level VARCHAR(20) DEFAULT 'free'
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_level ON subscriptions(level);
```

**Оценка времени:** 1 час

---

### 1.3. Настройка конфигурации (.env)

**Задача:** Создать конфигурационный файл для переменных окружения

**.env.example:**
```env
# Server Configuration
PORT=8080
ENV=development
FRONTEND_URL=http://localhost:3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=bull_run
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback

# Yandex OAuth
YANDEX_CLIENT_ID=your-yandex-client-id
YANDEX_CLIENT_SECRET=your-yandex-client-secret
YANDEX_REDIRECT_URL=http://localhost:8080/api/auth/yandex/callback

# Session Configuration
SESSION_COOKIE_NAME=bull_run_session
SESSION_MAX_AGE=86400
COOKIE_DOMAIN=localhost
COOKIE_SECURE=false

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1m
```

**Оценка времени:** 30 минут

---

## Фаза 2: Регистрация OAuth приложений

### 2.1. Регистрация Google OAuth приложения

**Задача:** Настроить OAuth 2.0 приложение в Google Cloud Console

**Шаги:**

1. Перейти в [Google Cloud Console](https://console.cloud.google.com/)
2. Создать новый проект "Bull Run Platform"
3. Включить Google+ API
4. Перейти в "APIs & Services" → "Credentials"
5. Создать "OAuth 2.0 Client ID"
6. Настроить OAuth consent screen:
   - App name: Bull Run
   - User support email: your-email@example.com
   - Developer contact: your-email@example.com
   - Scopes: email, profile, openid
7. Создать Web application credentials:
   - Name: Bull Run Backend
   - Authorized redirect URIs:
     - `http://localhost:8080/api/auth/google/callback` (dev)
     - `https://bull-run.com/api/auth/google/callback` (prod)
8. Сохранить Client ID и Client Secret в .env

**Необходимые OAuth Scopes:**
- `https://www.googleapis.com/auth/userinfo.email`
- `https://www.googleapis.com/auth/userinfo.profile`
- `openid`

**Получаемые данные от Google:**
```json
{
  "id": "123456789",
  "email": "user@gmail.com",
  "verified_email": true,
  "name": "John Doe",
  "given_name": "John",
  "family_name": "Doe",
  "picture": "https://lh3.googleusercontent.com/..."
}
```

**Оценка времени:** 1 час

---

### 2.2. Регистрация Яндекс OAuth приложения

**Задача:** Настроить OAuth приложение в Яндекс ID

**Шаги:**

1. Перейти в [Яндекс OAuth](https://oauth.yandex.ru/)
2. Нажать "Зарегистрировать новое приложение"
3. Заполнить данные:
   - Название: Bull Run
   - Описание: Платформа фундаментального анализа российских акций
   - Иконка: загрузить логотип
   - Права доступа:
     - Яндекс ID (login:email, login:info, login:avatar)
4. Указать Callback URI:
   - `http://localhost:8080/api/auth/yandex/callback` (dev)
   - `https://bull-run.com/api/auth/yandex/callback` (prod)
5. Сохранить Client ID и Client Secret в .env

**Необходимые OAuth Scopes:**
- `login:email` - доступ к email
- `login:info` - доступ к имени пользователя
- `login:avatar` - доступ к аватару

**Получаемые данные от Яндекс:**
```json
{
  "id": "123456789",
  "login": "user",
  "client_id": "abc123",
  "display_name": "John Doe",
  "real_name": "John Doe",
  "first_name": "John",
  "last_name": "Doe",
  "sex": "male",
  "default_email": "user@yandex.ru",
  "emails": ["user@yandex.ru"],
  "default_avatar_id": "123456/abc-def",
  "is_avatar_empty": false
}
```

**Оценка времени:** 1 час

---

## Фаза 3: Реализация бэкенда

### 3.1. Конфигурация и подключение к БД

**Задача:** Реализовать загрузку конфигурации и подключение к PostgreSQL/Redis

**internal/config/config.go:**
```go
package config

import (
    "os"
    "time"
    "github.com/joho/godotenv"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
    OAuth    OAuthConfig
}

type ServerConfig struct {
    Port        string
    Environment string
    FrontendURL string
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type RedisConfig struct {
    Host     string
    Port     string
    Password string
    DB       int
}

type JWTConfig struct {
    Secret         string
    Expiry         time.Duration
    RefreshExpiry  time.Duration
}

type OAuthConfig struct {
    Google GoogleOAuthConfig
    Yandex YandexOAuthConfig
}

type GoogleOAuthConfig struct {
    ClientID     string
    ClientSecret string
    RedirectURL  string
}

type YandexOAuthConfig struct {
    ClientID     string
    ClientSecret string
    RedirectURL  string
}

func Load() (*Config, error) {
    godotenv.Load()

    return &Config{
        Server: ServerConfig{
            Port:        getEnv("PORT", "8080"),
            Environment: getEnv("ENV", "development"),
            FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", ""),
            DBName:   getEnv("DB_NAME", "bull_run"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnv("REDIS_PORT", "6379"),
            Password: getEnv("REDIS_PASSWORD", ""),
            DB:       0,
        },
        JWT: JWTConfig{
            Secret:        getEnv("JWT_SECRET", ""),
            Expiry:        24 * time.Hour,
            RefreshExpiry: 7 * 24 * time.Hour,
        },
        OAuth: OAuthConfig{
            Google: GoogleOAuthConfig{
                ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
                ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
                RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
            },
            Yandex: YandexOAuthConfig{
                ClientID:     getEnv("YANDEX_CLIENT_ID", ""),
                ClientSecret: getEnv("YANDEX_CLIENT_SECRET", ""),
                RedirectURL:  getEnv("YANDEX_REDIRECT_URL", ""),
            },
        },
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

**Оценка времени:** 2 часа

---

### 3.2. Модели данных

**Задача:** Создать Go структуры для работы с БД

**internal/models/user.go:**
```go
package models

import "time"

type User struct {
    ID          int64      `json:"id" db:"id"`
    Name        string     `json:"name" db:"name"`
    LastLoginAt *time.Time `json:"lastLoginAt" db:"last_login_at"`
    CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
    UpdatedAt   *time.Time `json:"updatedAt" db:"updated_at"`
}

type Auth struct {
    UserID         int64  `json:"userId" db:"user_id"`
    Email          string `json:"email" db:"email"`
    HashedPassword string `json:"-" db:"hashed_password"`
}

type ProviderAuth struct {
    ID             int64      `json:"id" db:"id"`
    UserID         int64      `json:"userId" db:"user_id"`
    ProviderUserID string     `json:"providerUserId" db:"provider_user_id"`
    ProviderType   string     `json:"providerType" db:"provider_type"` // "google" or "yandex"
    Email          string     `json:"email" db:"email"`
    AvatarURL      *string    `json:"avatarUrl" db:"avatar_url"`
    AccessToken    *string    `json:"-" db:"access_token"`
    RefreshToken   *string    `json:"-" db:"refresh_token"`
    TokenExpiry    *time.Time `json:"-" db:"token_expiry"`
    CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
    UpdatedAt      *time.Time `json:"updatedAt" db:"updated_at"`
}

type Subscription struct {
    ID        int64      `json:"id" db:"id"`
    UserID    int64      `json:"userId" db:"user_id"`
    StartDate *time.Time `json:"startDate" db:"start_date"`
    EndDate   *time.Time `json:"endDate" db:"end_date"`
    Level     string     `json:"level" db:"level"` // "free", "premium", "pro"
}

// DTOs для API
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Name     string `json:"name" binding:"required,min=2"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    User         User   `json:"user"`
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
    ExpiresIn    int64  `json:"expiresIn"`
}

type OAuthUserInfo struct {
    ProviderType   string  `json:"providerType"`
    ProviderUserID string  `json:"providerUserId"`
    Email          string  `json:"email"`
    Name           string  `json:"name"`
    AvatarURL      *string `json:"avatarUrl"`
}
```

**Оценка времени:** 2 часа

---

### 3.3. JWT Service

**Задача:** Реализовать генерацию и верификацию JWT токенов

**internal/service/jwt_service.go:**
```go
package service

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "bull-run/internal/config"
)

type JWTService struct {
    config *config.Config
}

type Claims struct {
    UserID int64  `json:"userId"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func NewJWTService(cfg *config.Config) *JWTService {
    return &JWTService{config: cfg}
}

func (s *JWTService) GenerateAccessToken(userID int64, email string) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.Expiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "bull-run",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *JWTService) GenerateRefreshToken(userID int64, email string) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.RefreshExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "bull-run",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(s.config.JWT.Secret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}
```

**Оценка времени:** 2 часа

---

### 3.4. OAuth Service

**Задача:** Реализовать OAuth flow для Google и Яндекс

**internal/service/oauth_service.go:**
```go
package service

import (
    "context"
    "encoding/json"
    "errors"
    "io"
    "net/http"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "golang.org/x/oauth2/yandex"

    "bull-run/internal/config"
    "bull-run/internal/models"
)

type OAuthService struct {
    config       *config.Config
    googleConfig *oauth2.Config
    yandexConfig *oauth2.Config
}

func NewOAuthService(cfg *config.Config) *OAuthService {
    return &OAuthService{
        config: cfg,
        googleConfig: &oauth2.Config{
            ClientID:     cfg.OAuth.Google.ClientID,
            ClientSecret: cfg.OAuth.Google.ClientSecret,
            RedirectURL:  cfg.OAuth.Google.RedirectURL,
            Scopes: []string{
                "https://www.googleapis.com/auth/userinfo.email",
                "https://www.googleapis.com/auth/userinfo.profile",
            },
            Endpoint: google.Endpoint,
        },
        yandexConfig: &oauth2.Config{
            ClientID:     cfg.OAuth.Yandex.ClientID,
            ClientSecret: cfg.OAuth.Yandex.ClientSecret,
            RedirectURL:  cfg.OAuth.Yandex.RedirectURL,
            Scopes:       []string{"login:email", "login:info", "login:avatar"},
            Endpoint:     yandex.Endpoint,
        },
    }
}

// Google OAuth
func (s *OAuthService) GetGoogleAuthURL(state string) string {
    return s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *OAuthService) ExchangeGoogleCode(ctx context.Context, code string) (*oauth2.Token, error) {
    return s.googleConfig.Exchange(ctx, code)
}

func (s *OAuthService) GetGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*models.OAuthUserInfo, error) {
    client := s.googleConfig.Client(ctx, token)

    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var googleUser struct {
        ID      string `json:"id"`
        Email   string `json:"email"`
        Name    string `json:"name"`
        Picture string `json:"picture"`
    }

    if err := json.Unmarshal(data, &googleUser); err != nil {
        return nil, err
    }

    return &models.OAuthUserInfo{
        ProviderType:   "google",
        ProviderUserID: googleUser.ID,
        Email:          googleUser.Email,
        Name:           googleUser.Name,
        AvatarURL:      &googleUser.Picture,
    }, nil
}

// Yandex OAuth
func (s *OAuthService) GetYandexAuthURL(state string) string {
    return s.yandexConfig.AuthCodeURL(state)
}

func (s *OAuthService) ExchangeYandexCode(ctx context.Context, code string) (*oauth2.Token, error) {
    return s.yandexConfig.Exchange(ctx, code)
}

func (s *OAuthService) GetYandexUserInfo(ctx context.Context, token *oauth2.Token) (*models.OAuthUserInfo, error) {
    client := s.yandexConfig.Client(ctx, token)

    resp, err := client.Get("https://login.yandex.ru/info")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var yandexUser struct {
        ID           string `json:"id"`
        DisplayName  string `json:"display_name"`
        DefaultEmail string `json:"default_email"`
        AvatarID     string `json:"default_avatar_id"`
        IsAvatarEmpty bool  `json:"is_avatar_empty"`
    }

    if err := json.Unmarshal(data, &yandexUser); err != nil {
        return nil, err
    }

    var avatarURL *string
    if !yandexUser.IsAvatarEmpty && yandexUser.AvatarID != "" {
        url := "https://avatars.yandex.net/get-yapic/" + yandexUser.AvatarID + "/islands-200"
        avatarURL = &url
    }

    return &models.OAuthUserInfo{
        ProviderType:   "yandex",
        ProviderUserID: yandexUser.ID,
        Email:          yandexUser.DefaultEmail,
        Name:           yandexUser.DisplayName,
        AvatarURL:      avatarURL,
    }, nil
}
```

**Оценка времени:** 4 часа

---

### 3.5. Repository Layer

**Задача:** Реализовать слой работы с БД

**internal/repository/user_repository.go:**
```go
package repository

import (
    "database/sql"
    "time"

    "bull-run/internal/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(name string) (*models.User, error) {
    user := &models.User{}
    now := time.Now()

    err := r.db.QueryRow(
        `INSERT INTO users (name, created_at) VALUES ($1, $2) RETURNING id, name, created_at`,
        name, now,
    ).Scan(&user.ID, &user.Name, &user.CreatedAt)

    return user, err
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
    user := &models.User{}
    err := r.db.QueryRow(
        `SELECT id, name, last_login_at, created_at, updated_at FROM users WHERE id = $1`,
        id,
    ).Scan(&user.ID, &user.Name, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)

    if err == sql.ErrNoRows {
        return nil, nil
    }
    return user, err
}

func (r *UserRepository) UpdateLastLogin(id int64) error {
    _, err := r.db.Exec(
        `UPDATE users SET last_login_at = $1 WHERE id = $2`,
        time.Now(), id,
    )
    return err
}
```

**internal/repository/provider_repository.go:**
```go
package repository

import (
    "database/sql"
    "time"

    "bull-run/internal/models"
)

type ProviderRepository struct {
    db *sql.DB
}

func NewProviderRepository(db *sql.DB) *ProviderRepository {
    return &ProviderRepository{db: db}
}

func (r *ProviderRepository) FindByProvider(providerType, providerUserID string) (*models.ProviderAuth, error) {
    auth := &models.ProviderAuth{}
    err := r.db.QueryRow(
        `SELECT id, user_id, provider_user_id, provider_type, email, avatar_url, created_at, updated_at
         FROM provider_auth WHERE provider_type = $1 AND provider_user_id = $2`,
        providerType, providerUserID,
    ).Scan(
        &auth.ID, &auth.UserID, &auth.ProviderUserID, &auth.ProviderType,
        &auth.Email, &auth.AvatarURL, &auth.CreatedAt, &auth.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    return auth, err
}

func (r *ProviderRepository) Create(userID int64, info *models.OAuthUserInfo) error {
    now := time.Now()
    _, err := r.db.Exec(
        `INSERT INTO provider_auth
         (user_id, provider_user_id, provider_type, email, avatar_url, created_at)
         VALUES ($1, $2, $3, $4, $5, $6)`,
        userID, info.ProviderUserID, info.ProviderType, info.Email, info.AvatarURL, now,
    )
    return err
}

func (r *ProviderRepository) UpdateTokens(id int64, accessToken, refreshToken string, expiry time.Time) error {
    _, err := r.db.Exec(
        `UPDATE provider_auth
         SET access_token = $1, refresh_token = $2, token_expiry = $3, updated_at = $4
         WHERE id = $5`,
        accessToken, refreshToken, expiry, time.Now(), id,
    )
    return err
}
```

**Оценка времени:** 3 часа

---

### 3.6. Auth Service (бизнес-логика)

**Задача:** Реализовать основную логику аутентификации

**internal/service/auth_service.go:**
```go
package service

import (
    "context"
    "database/sql"
    "errors"

    "bull-run/internal/models"
    "bull-run/internal/repository"
)

type AuthService struct {
    userRepo     *repository.UserRepository
    providerRepo *repository.ProviderRepository
    jwtService   *JWTService
    oauthService *OAuthService
}

func NewAuthService(
    userRepo *repository.UserRepository,
    providerRepo *repository.ProviderRepository,
    jwtService *JWTService,
    oauthService *OAuthService,
) *AuthService {
    return &AuthService{
        userRepo:     userRepo,
        providerRepo: providerRepo,
        jwtService:   jwtService,
        oauthService: oauthService,
    }
}

func (s *AuthService) HandleOAuthCallback(
    ctx context.Context,
    providerType string,
    code string,
) (*models.AuthResponse, error) {
    var userInfo *models.OAuthUserInfo
    var err error

    // Получаем токен и информацию о пользователе
    switch providerType {
    case "google":
        token, err := s.oauthService.ExchangeGoogleCode(ctx, code)
        if err != nil {
            return nil, err
        }
        userInfo, err = s.oauthService.GetGoogleUserInfo(ctx, token)
        if err != nil {
            return nil, err
        }
    case "yandex":
        token, err := s.oauthService.ExchangeYandexCode(ctx, code)
        if err != nil {
            return nil, err
        }
        userInfo, err = s.oauthService.GetYandexUserInfo(ctx, token)
        if err != nil {
            return nil, err
        }
    default:
        return nil, errors.New("unsupported provider")
    }

    // Ищем существующую связку
    providerAuth, err := s.providerRepo.FindByProvider(userInfo.ProviderType, userInfo.ProviderUserID)
    if err != nil {
        return nil, err
    }

    var user *models.User

    if providerAuth != nil {
        // Пользователь уже существует
        user, err = s.userRepo.GetByID(providerAuth.UserID)
        if err != nil {
            return nil, err
        }
    } else {
        // Создаем нового пользователя
        user, err = s.userRepo.Create(userInfo.Name)
        if err != nil {
            return nil, err
        }

        // Создаем связку с OAuth провайдером
        err = s.providerRepo.Create(user.ID, userInfo)
        if err != nil {
            return nil, err
        }

        // Создаем бесплатную подписку
        // TODO: implement subscription creation
    }

    // Обновляем время последнего входа
    s.userRepo.UpdateLastLogin(user.ID)

    // Генерируем JWT токены
    accessToken, err := s.jwtService.GenerateAccessToken(user.ID, userInfo.Email)
    if err != nil {
        return nil, err
    }

    refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, userInfo.Email)
    if err != nil {
        return nil, err
    }

    return &models.AuthResponse{
        User:         *user,
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    86400, // 24 hours
    }, nil
}
```

**Оценка времени:** 3 часа

---

### 3.7. HTTP Handlers

**Задача:** Реализовать HTTP обработчики для OAuth endpoints

**internal/handler/oauth_handler.go:**
```go
package handler

import (
    "net/http"
    "crypto/rand"
    "encoding/base64"

    "github.com/gin-gonic/gin"
    "bull-run/internal/service"
    "bull-run/internal/config"
)

type OAuthHandler struct {
    authService  *service.AuthService
    oauthService *service.OAuthService
    config       *config.Config
}

func NewOAuthHandler(
    authService *service.AuthService,
    oauthService *service.OAuthService,
    cfg *config.Config,
) *OAuthHandler {
    return &OAuthHandler{
        authService:  authService,
        oauthService: oauthService,
        config:       cfg,
    }
}

// Google OAuth
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
    state := generateState()

    // Сохраняем state в cookie для верификации
    c.SetCookie("oauth_state", state, 300, "/", "", false, true)

    url := h.oauthService.GetGoogleAuthURL(state)
    c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
    // Проверяем state
    savedState, err := c.Cookie("oauth_state")
    if err != nil || savedState != c.Query("state") {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
        return
    }

    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
        return
    }

    authResponse, err := h.authService.HandleOAuthCallback(c.Request.Context(), "google", code)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Устанавливаем токены в cookies
    h.setAuthCookies(c, authResponse.AccessToken, authResponse.RefreshToken)

    // Редирект на фронтенд
    c.Redirect(http.StatusTemporaryRedirect, h.config.Server.FrontendURL+"/dashboard")
}

// Yandex OAuth
func (h *OAuthHandler) YandexLogin(c *gin.Context) {
    state := generateState()

    c.SetCookie("oauth_state", state, 300, "/", "", false, true)

    url := h.oauthService.GetYandexAuthURL(state)
    c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *OAuthHandler) YandexCallback(c *gin.Context) {
    savedState, err := c.Cookie("oauth_state")
    if err != nil || savedState != c.Query("state") {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
        return
    }

    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
        return
    }

    authResponse, err := h.authService.HandleOAuthCallback(c.Request.Context(), "yandex", code)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    h.setAuthCookies(c, authResponse.AccessToken, authResponse.RefreshToken)

    c.Redirect(http.StatusTemporaryRedirect, h.config.Server.FrontendURL+"/dashboard")
}

// Вспомогательные функции
func (h *OAuthHandler) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
    c.SetCookie("access_token", accessToken, 86400, "/", "", false, true)
    c.SetCookie("refresh_token", refreshToken, 604800, "/", "", false, true)
}

func generateState() string {
    b := make([]byte, 16)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
```

**Оценка времени:** 3 часа

---

### 3.8. Middleware для авторизации

**Задача:** Создать middleware для проверки JWT токенов

**internal/middleware/auth_middleware.go:**
```go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "bull-run/internal/service"
)

func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Пробуем получить токен из заголовка
        authHeader := c.GetHeader("Authorization")
        var tokenString string

        if authHeader != "" {
            parts := strings.Split(authHeader, " ")
            if len(parts) == 2 && parts[0] == "Bearer" {
                tokenString = parts[1]
            }
        }

        // Если в заголовке нет, пробуем из cookie
        if tokenString == "" {
            var err error
            tokenString, err = c.Cookie("access_token")
            if err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
                c.Abort()
                return
            }
        }

        // Валидируем токен
        claims, err := jwtService.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Сохраняем данные пользователя в контексте
        c.Set("userId", claims.UserID)
        c.Set("email", claims.Email)

        c.Next()
    }
}
```

**Оценка времени:** 1 час

---

### 3.9. Роутинг и главный файл

**Задача:** Настроить маршруты и точку входа

**cmd/api/main.go:**
```go
package main

import (
    "database/sql"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    _ "github.com/lib/pq"

    "bull-run/internal/config"
    "bull-run/internal/handler"
    "bull-run/internal/middleware"
    "bull-run/internal/repository"
    "bull-run/internal/service"
)

func main() {
    // Загружаем конфигурацию
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Подключаемся к БД
    db, err := sql.Open("postgres", buildDSN(cfg))
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Инициализируем репозитории
    userRepo := repository.NewUserRepository(db)
    providerRepo := repository.NewProviderRepository(db)

    // Инициализируем сервисы
    jwtService := service.NewJWTService(cfg)
    oauthService := service.NewOAuthService(cfg)
    authService := service.NewAuthService(userRepo, providerRepo, jwtService, oauthService)

    // Инициализируем handlers
    oauthHandler := handler.NewOAuthHandler(authService, oauthService, cfg)

    // Настраиваем Gin
    r := gin.Default()

    // CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{cfg.Server.FrontendURL},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))

    // Роуты
    api := r.Group("/api")
    {
        auth := api.Group("/auth")
        {
            // Google OAuth
            auth.GET("/google", oauthHandler.GoogleLogin)
            auth.GET("/google/callback", oauthHandler.GoogleCallback)

            // Yandex OAuth
            auth.GET("/yandex", oauthHandler.YandexLogin)
            auth.GET("/yandex/callback", oauthHandler.YandexCallback)

            // Logout
            auth.POST("/logout", func(c *gin.Context) {
                c.SetCookie("access_token", "", -1, "/", "", false, true)
                c.SetCookie("refresh_token", "", -1, "/", "", false, true)
                c.JSON(200, gin.H{"message": "Logged out"})
            })
        }

        // Защищенные роуты
        protected := api.Group("")
        protected.Use(middleware.AuthMiddleware(jwtService))
        {
            protected.GET("/me", func(c *gin.Context) {
                userID := c.GetInt64("userId")
                user, err := userRepo.GetByID(userID)
                if err != nil {
                    c.JSON(500, gin.H{"error": "Failed to fetch user"})
                    return
                }
                c.JSON(200, user)
            })
        }
    }

    // Запуск сервера
    log.Printf("Server starting on port %s", cfg.Server.Port)
    r.Run(":" + cfg.Server.Port)
}

func buildDSN(cfg *config.Config) string {
    return "host=" + cfg.Database.Host +
        " port=" + cfg.Database.Port +
        " user=" + cfg.Database.User +
        " password=" + cfg.Database.Password +
        " dbname=" + cfg.Database.DBName +
        " sslmode=" + cfg.Database.SSLMode
}
```

**Оценка времени:** 2 часа

---

## Фаза 4: Интеграция с фронтендом

### 4.1. Обновление API клиента

**Задача:** Настроить фронтенд для работы с OAuth endpoints

**frontend/lib/api.ts:** (обновить существующий файл)
```typescript
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'

export const authAPI = {
  // Google OAuth
  async loginWithGoogle() {
    const response = await fetch(`${API_BASE_URL}/auth/google`, {
      credentials: 'include',
    })
    const data = await response.json()
    window.location.href = data.url
  },

  // Yandex OAuth
  async loginWithYandex() {
    const response = await fetch(`${API_BASE_URL}/auth/yandex`, {
      credentials: 'include',
    })
    const data = await response.json()
    window.location.href = data.url
  },

  // Get current user
  async getCurrentUser() {
    const response = await fetch(`${API_BASE_URL}/me`, {
      credentials: 'include',
    })
    if (!response.ok) throw new Error('Not authenticated')
    return response.json()
  },

  // Logout
  async logout() {
    await fetch(`${API_BASE_URL}/auth/logout`, {
      method: 'POST',
      credentials: 'include',
    })
  },
}
```

**Оценка времени:** 1 час

---

### 4.2. Обновление страницы авторизации

**Задача:** Подключить реальные OAuth функции к кнопкам

**frontend/app/auth/page.tsx:** (обновить обработчики)
```typescript
'use client'

import { authAPI } from '@/lib/api'

export default function AuthPage() {
  const handleGoogleLogin = async () => {
    try {
      await authAPI.loginWithGoogle()
    } catch (error) {
      console.error('Google login failed:', error)
    }
  }

  const handleYandexLogin = async () => {
    try {
      await authAPI.loginWithYandex()
    } catch (error) {
      console.error('Yandex login failed:', error)
    }
  }

  // Остальной код остается прежним, только добавляем onClick handlers
}
```

**Оценка времени:** 30 минут

---

### 4.3. Создание AuthContext

**Задача:** Создать контекст для управления состоянием авторизации

**frontend/contexts/AuthContext.tsx:**
```typescript
'use client'

import { createContext, useContext, useEffect, useState } from 'react'
import { authAPI } from '@/lib/api'
import { User } from '@/types/user'

interface AuthContextType {
  user: User | null
  loading: boolean
  logout: () => Promise<void>
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    checkAuth()
  }, [])

  const checkAuth = async () => {
    try {
      const userData = await authAPI.getCurrentUser()
      setUser(userData)
    } catch (error) {
      setUser(null)
    } finally {
      setLoading(false)
    }
  }

  const logout = async () => {
    await authAPI.logout()
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, loading, logout, refreshUser: checkAuth }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within AuthProvider')
  return context
}
```

**Оценка времени:** 1 час

---

## Фаза 5: Тестирование и деплой

### 5.1. Локальное тестирование

**Задача:** Протестировать OAuth flow локально

**Чек-лист:**
- [ ] Google OAuth login работает
- [ ] Yandex OAuth login работает
- [ ] Токены сохраняются в cookies
- [ ] Пользователь создается в БД
- [ ] Subscription создается с уровнем "free"
- [ ] Редирект на dashboard после авторизации
- [ ] Logout очищает cookies
- [ ] Protected endpoints требуют авторизацию
- [ ] Refresh token работает
- [ ] CORS настроен правильно

**Оценка времени:** 4 часа

---

### 5.2. Настройка Docker

**Задача:** Добавить бэкенд в Docker Compose

**compose.yaml:** (обновить)
```yaml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: bull_run
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    env_file:
      - ./backend/.env
    depends_on:
      - postgres
      - redis
    volumes:
      - ./backend:/app

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080/api
    depends_on:
      - backend

volumes:
  postgres_data:
```

**backend/Dockerfile:**
```dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bull-run-api ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /bull-run-api .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./bull-run-api"]
```

**Оценка времени:** 2 часа

---

### 5.3. Production deployment

**Задача:** Подготовить к production деплою

**Чек-лист:**
- [ ] Обновить Callback URLs в Google Cloud Console (prod URLs)
- [ ] Обновить Callback URLs в Yandex OAuth (prod URLs)
- [ ] Настроить HTTPS
- [ ] Настроить COOKIE_SECURE=true
- [ ] Настроить CORS для production домена
- [ ] Добавить rate limiting
- [ ] Настроить мониторинг и логирование
- [ ] Создать backup стратегию для БД
- [ ] Настроить CI/CD pipeline

**Оценка времени:** 6 часов

---

## Итоговая оценка времени

| Фаза | Задачи | Время |
|------|--------|-------|
| **Фаза 1: Подготовка** | Структура, миграции, конфиг | 5-6 часов |
| **Фаза 2: OAuth Apps** | Google и Яндекс регистрация | 2 часа |
| **Фаза 3: Backend** | Go сервисы, handlers, middleware | 20-22 часа |
| **Фаза 4: Frontend** | Интеграция с UI | 2-3 часа |
| **Фаза 5: Testing & Deploy** | Тестирование, Docker, деплой | 12 часов |
| **Итого** | | **41-45 часов** |

---

## Безопасность

### Критические моменты:

1. **CSRF Protection:**
   - Используем `state` параметр в OAuth flow
   - Проверяем state при callback

2. **XSS Protection:**
   - HTTPOnly cookies для токенов
   - Валидация всех входящих данных

3. **Секреты:**
   - Никогда не коммитим `.env` файл
   - Используем переменные окружения в production

4. **JWT Security:**
   - Используем strong secret key (256 бит)
   - Короткий expiry для access tokens (24h)
   - Refresh tokens для продления сессии

5. **Rate Limiting:**
   - Ограничить количество OAuth попыток
   - Защита от brute-force атак

6. **HTTPS:**
   - Обязательно в production
   - Secure cookies только через HTTPS

---

## Дополнительные улучшения (опционально)

1. **Email + Password авторизация:**
   - Хэширование паролей (bcrypt)
   - Email верификация
   - Восстановление пароля

2. **Множественные OAuth провайдеры:**
   - Возможность привязать несколько провайдеров к одному аккаунту

3. **2FA (Two-Factor Authentication):**
   - TOTP (Google Authenticator)
   - SMS коды

4. **Session Management:**
   - Хранение активных сессий в Redis
   - Возможность logout со всех устройств

5. **Audit Log:**
   - Логирование всех входов/выходов
   - История авторизаций

---

## Финальный чек-лист перед запуском

- [ ] Все миграции БД применены
- [ ] Google OAuth credentials настроены
- [ ] Yandex OAuth credentials настроены
- [ ] Переменные окружения заполнены
- [ ] JWT secret сгенерирован (сильный)
- [ ] CORS настроен для фронтенда
- [ ] Cookies работают (HTTPOnly, Secure в prod)
- [ ] OAuth state верификация работает
- [ ] Пользователи создаются в БД
- [ ] Подписки создаются автоматически
- [ ] Редиректы работают корректно
- [ ] Защищенные эндпоинты требуют auth
- [ ] Logout очищает сессию
- [ ] Error handling реализован
- [ ] Логирование настроено
- [ ] Тесты пройдены

---

## Документация и ресурсы

**OAuth 2.0 документация:**
- [Google OAuth 2.0](https://developers.google.com/identity/protocols/oauth2)
- [Яндекс OAuth](https://yandex.ru/dev/id/doc/ru/)

**Go библиотеки:**
- [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)
- [github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt)

**Best Practices:**
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [OAuth 2.0 Security Best Practices](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics)
