package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth-service/internal/domain"
)

const defaultRefreshTokenTTL = 15 * 24 * time.Hour

type DbRepo struct {
	pool            *pgxpool.Pool
	refreshTokenTTL time.Duration
}

func NewDbRepo(pool *pgxpool.Pool) *DbRepo {
	return &DbRepo{
		pool:            pool,
		refreshTokenTTL: defaultRefreshTokenTTL,
	}
}

func NewDbRepoWithTTL(pool *pgxpool.Pool, refreshTokenTTL time.Duration) *DbRepo {
	return &DbRepo{
		pool:            pool,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (r *DbRepo) CreateUser(ctx context.Context, name, email, hashedPassword string) (*domain.User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var userID int64
	var createdAt, updatedAt time.Time
	var lastLoginAt sql.NullTime

	err = tx.QueryRow(ctx, `
		INSERT INTO users (name, status, last_login_at, created_at, updated_at)
		VALUES ($1, 'active', NULL, NOW(), NOW())
		RETURNING id, created_at, updated_at, last_login_at
	`, name).Scan(&userID, &createdAt, &updatedAt, &lastLoginAt)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO auth (user_id, email, password_hash)
		VALUES ($1, $2, $3)
	`, userID, email, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("insert auth data: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	var lastLogin time.Time
	if lastLoginAt.Valid {
		lastLogin = lastLoginAt.Time
	}

	return domain.NewUser(userID, name, "active", lastLogin, createdAt, updatedAt), nil
}

func (r *DbRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var id int64
	var name, status string
	var lastLoginAt sql.NullTime
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx, `
		SELECT u.id, u.name, u.status, u.last_login_at, u.created_at, u.updated_at
		FROM users u
		INNER JOIN auth a ON u.id = a.user_id
		WHERE a.email = $1
	`, email).Scan(
		&id,
		&name,
		&status,
		&lastLoginAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	var lastLogin time.Time
	if lastLoginAt.Valid {
		lastLogin = lastLoginAt.Time
	}

	return domain.NewUser(id, name, status, lastLogin, createdAt, updatedAt), nil
}

func (r *DbRepo) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	var userID int64
	var name, status string
	var lastLoginAt sql.NullTime
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx, `
		SELECT id, name, status, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&userID,
		&name,
		&status,
		&lastLoginAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	var lastLogin time.Time
	if lastLoginAt.Valid {
		lastLogin = lastLoginAt.Time
	}

	return domain.NewUser(userID, name, status, lastLogin, createdAt, updatedAt), nil
}

func (r *DbRepo) UpdateLastLogin(ctx context.Context, userID int64) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users
		SET last_login_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, userID)
	if err != nil {
		return fmt.Errorf("update last login: %w", err)
	}

	return nil
}

func (r *DbRepo) CreateAuthData(ctx context.Context, userID int64, email, passwordHash string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO auth (user_id, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
		SET email = EXCLUDED.email, password_hash = EXCLUDED.password_hash
	`, userID, email, passwordHash)
	if err != nil {
		return fmt.Errorf("create auth data: %w", err)
	}

	return nil
}

func (r *DbRepo) GetAuthDataByEmail(ctx context.Context, email string) (*domain.AuthData, error) {
	var userID int64
	var emailStr, passwordHash string

	err := r.pool.QueryRow(ctx, `
		SELECT user_id, email, password_hash
		FROM auth
		WHERE email = $1
	`, email).Scan(
		&userID,
		&emailStr,
		&passwordHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get auth data by email: %w", err)
	}

	return domain.NewAuthData(userID, emailStr, passwordHash), nil
}

func (r *DbRepo) CreateProviderAuth(ctx context.Context, userID int64, providerType domain.ProviderType, providerUserID string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO provider_auth (user_id, provider_type, provider_user_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (provider_type, provider_user_id) DO NOTHING
	`, userID, string(providerType), providerUserID)
	if err != nil {
		return fmt.Errorf("create provider auth: %w", err)
	}

	return nil
}

func (r *DbRepo) GetProviderAuth(ctx context.Context, providerType domain.ProviderType, providerUserID string) (*domain.ProviderAuth, error) {
	var id, userID int64
	var providerUserIDStr string
	var providerTypeStr string
	var email sql.NullString

	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, provider_user_id, provider_type, email
		FROM provider_auth
		WHERE provider_type = $1 AND provider_user_id = $2
	`, string(providerType), providerUserID).Scan(
		&id,
		&userID,
		&providerUserIDStr,
		&providerTypeStr,
		&email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get provider auth: %w", err)
	}

	emailStr := ""
	if email.Valid {
		emailStr = email.String
	}

	return domain.NewProviderAuth(id, userID, providerUserIDStr, domain.ProviderType(providerTypeStr), emailStr), nil
}

func (r *DbRepo) GetProvidersByUserID(ctx context.Context, userID int64) ([]domain.ProviderAuth, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, provider_user_id, provider_type, email
		FROM provider_auth
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("get providers by user id: %w", err)
	}
	defer rows.Close()

	var providers []domain.ProviderAuth
	for rows.Next() {
		var id, userID int64
		var providerUserIDStr, providerTypeStr string
		var email sql.NullString

		if err := rows.Scan(
			&id,
			&userID,
			&providerUserIDStr,
			&providerTypeStr,
			&email,
		); err != nil {
			return nil, fmt.Errorf("scan provider auth: %w", err)
		}

		emailStr := ""
		if email.Valid {
			emailStr = email.String
		}

		provider := domain.NewProviderAuth(id, userID, providerUserIDStr, domain.ProviderType(providerTypeStr), emailStr)
		providers = append(providers, *provider)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate providers: %w", err)
	}

	return providers, nil
}

func (r *DbRepo) CreateRefreshToken(ctx context.Context, tokenHash string, userID int64, deviceInfo string) error {
	expiresAt := time.Now().Add(r.refreshTokenTTL)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, device_info, created_at, updated_at, expires_at)
		VALUES ($1, $2, $3, NOW(), NOW(), $4)
	`, userID, tokenHash, deviceInfo, expiresAt)
	if err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}

	return nil
}

func (r *DbRepo) GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	var id, userID int64
	var tokenHashStr, deviceInfo string
	var createdAt, updatedAt, expiresAt time.Time

	err := r.pool.QueryRow(ctx, `
		SELECT id, token_hash, user_id, device_info, created_at, updated_at, expires_at
		FROM refresh_tokens
		WHERE token_hash = $1 AND expires_at > NOW()
	`, tokenHash).Scan(
		&id,
		&tokenHashStr,
		&userID,
		&deviceInfo,
		&createdAt,
		&updatedAt,
		&expiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	return domain.NewRefreshToken(id, tokenHashStr, userID, deviceInfo, createdAt, updatedAt, expiresAt), nil
}

func (r *DbRepo) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE refresh_tokens
		SET expires_at = NOW() - INTERVAL '1 second', updated_at = NOW()
		WHERE token_hash = $1
	`, tokenHash)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}

func (r *DbRepo) RevokeAllUserTokens(ctx context.Context, userID int64) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE refresh_tokens
		SET expires_at = NOW() - INTERVAL '1 second', updated_at = NOW()
		WHERE user_id = $1 AND expires_at > NOW()
	`, userID)
	if err != nil {
		return fmt.Errorf("revoke all user tokens: %w", err)
	}

	return nil
}
