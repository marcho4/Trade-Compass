package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID int64  `json:"sub"`
	Name   string `json:"name"`
	Status string `json:"status"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID  int64  `json:"sub"`
	TokenID string `json:"jti"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type JWTConfig struct {
	SecretKey       []byte
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type JWTService interface {
	GenerateAccessToken(userID int64, name, status string) (string, error)
	GenerateRefreshToken(userID int64) (token string, tokenID string, err error)
	ValidateAccessToken(token string) (*AccessClaims, error)
	ValidateRefreshToken(token string) (*RefreshClaims, error)
}
