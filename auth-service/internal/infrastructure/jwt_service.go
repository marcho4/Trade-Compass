package infrastructure

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"auth-service/internal/domain"
)

type jwtService struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTService(cfg domain.JWTConfig) (domain.JWTService, error) {
	if len(cfg.SecretKey) < 32 {
		return nil, fmt.Errorf("secret key must be at least 32 bytes")
	}

	return &jwtService{
		secretKey:       cfg.SecretKey,
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}, nil
}

func (s *jwtService) GenerateAccessToken(userID int64, name, status string) (string, error) {
	now := time.Now()

	claims := domain.AccessClaims{
		UserID: userID,
		Name:   name,
		Status: status,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signedToken, nil
}

func (s *jwtService) GenerateRefreshToken(userID int64) (string, string, error) {
	now := time.Now()
	tokenID := uuid.New().String()

	claims := domain.RefreshClaims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", "", fmt.Errorf("sign refresh token: %w", err)
	}

	return signedToken, tokenID, nil
}

func (s *jwtService) ValidateAccessToken(tokenString string) (*domain.AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse access token: %w", err)
	}

	claims, ok := token.Claims.(*domain.AccessClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token claims")
	}

	return claims, nil
}

func (s *jwtService) ValidateRefreshToken(tokenString string) (*domain.RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(*domain.RefreshClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	return claims, nil
}
