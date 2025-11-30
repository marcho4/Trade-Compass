package domain

import "context"

type DbRepo interface {
	CreateUser(ctx context.Context, name, email, hashedPassword string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uint64) (*User, error)
	UpdateLastLogin(ctx context.Context, userID int64) error

	CreateAuthData(ctx context.Context, userID int64, email, passwordHash string) error
	GetAuthDataByEmail(ctx context.Context, email string) (*AuthData, error)

	CreateProviderAuth(ctx context.Context, userID int64, providerType ProviderType, providerUserID string) error
	GetProviderAuth(ctx context.Context, providerType ProviderType, providerUserID string) (*ProviderAuth, error)
	GetProvidersByUserID(ctx context.Context, userID int64) ([]ProviderAuth, error)

	CreateRefreshToken(ctx context.Context, tokenHash string, userID int64, deviceInfo string) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID int64) error
}
