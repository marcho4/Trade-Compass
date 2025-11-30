package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"auth-service/internal/domain"
)

type Service struct {
	dbRepo             domain.DbRepo
	jwtService         domain.JWTService
	yandexClientId     string
	yandexClientSecret string
}

func NewService(dbRepo domain.DbRepo, jwtService domain.JWTService) *Service {
	yandexClientId := os.Getenv("YANDEX_CLIENT_ID")
	yandexClientSecret := os.Getenv("YANDEX_CLIENT_SECRET")

	if yandexClientSecret == "" || yandexClientId == "" {
		log.Fatal("Yandex credentials are not set")
	}

	return &Service{
		dbRepo:             dbRepo,
		jwtService:         jwtService,
		yandexClientId:     yandexClientId,
		yandexClientSecret: yandexClientSecret,
	}
}

func (s *Service) CreateUser(ctx context.Context, name, email, hashedPassword, deviceInfo string) (string, string, error) {
	user, err := s.dbRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("get user by email: %w", err)
	}

	if user != nil {
		return "", "", errors.New("user already exists")
	}

	user, err = s.dbRepo.CreateUser(ctx, name, email, hashedPassword)
	if err != nil {
		return "", "", fmt.Errorf("create user: %w", err)
	}

	accessToken, refreshToken, err := s.GenerateTokens(ctx, user.GetID(), user.GetName(), user.GetStatus(), deviceInfo)
	if err != nil {
		return "", "", fmt.Errorf("generate tokens: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *Service) ProcessLogin(ctx context.Context, email, hashedPassword, deviceInfo string) (string, string, error) {
	authData, err := s.dbRepo.GetAuthDataByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("get auth data: %w", err)
	}

	if authData == nil || authData.GetPasswordHash() != hashedPassword || authData.GetEmail() != email {
		return "", "", errors.New("invalid credentials")
	}

	user, err := s.dbRepo.GetUserByID(ctx, uint64(authData.GetUserID()))
	if err != nil {
		return "", "", fmt.Errorf("get user: %w", err)
	}

	accessToken, refreshToken, err := s.GenerateTokens(ctx, authData.GetUserID(), user.GetName(), user.GetStatus(), deviceInfo)
	if err != nil {
		return "", "", fmt.Errorf("generate tokens: %w", err)
	}

	if err := s.dbRepo.UpdateLastLogin(ctx, authData.GetUserID()); err != nil {
		return "", "", fmt.Errorf("update last login: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *Service) GenerateTokens(ctx context.Context, userID int64, name, status, deviceInfo string) (string, string, error) {
	accessToken, err := s.jwtService.GenerateAccessToken(userID, name, status)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, tokenID, err := s.jwtService.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	if err := s.dbRepo.CreateRefreshToken(ctx, tokenID, userID, deviceInfo); err != nil {
		return "", "", fmt.Errorf("save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *Service) RefreshTokens(ctx context.Context, refreshTokenString, deviceInfo string) (string, string, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", fmt.Errorf("validate refresh token: %w", err)
	}

	storedToken, err := s.dbRepo.GetRefreshToken(ctx, claims.TokenID)
	if err != nil {
		return "", "", fmt.Errorf("get stored refresh token: %w", err)
	}
	if storedToken == nil {
		return "", "", errors.New("refresh token revoked or not found")
	}

	if err := s.dbRepo.RevokeRefreshToken(ctx, claims.TokenID); err != nil {
		return "", "", fmt.Errorf("revoke old refresh token: %w", err)
	}

	user, err := s.dbRepo.GetUserByID(ctx, uint64(claims.UserID))
	if err != nil {
		return "", "", fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return "", "", errors.New("user not found")
	}

	return s.GenerateTokens(ctx, claims.UserID, user.GetName(), user.GetStatus(), deviceInfo)
}

func (s *Service) RevokeRefreshToken(ctx context.Context, refreshTokenString string) error {
	claims, err := s.jwtService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return fmt.Errorf("validate refresh token: %w", err)
	}

	if err := s.dbRepo.RevokeRefreshToken(ctx, claims.TokenID); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}

func (s *Service) RevokeAllUserTokens(ctx context.Context, userID int64) error {
	if err := s.dbRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		return fmt.Errorf("revoke all user tokens: %w", err)
	}

	return nil
}

func (s *Service) ValidateAccessToken(token string) (*domain.AccessClaims, error) {
	return s.jwtService.ValidateAccessToken(token)
}

func (s *Service) GetYandexAccessToken(yandexToken string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", yandexToken)
	data.Set("client_id", s.yandexClientId)
	data.Set("client_secret", s.yandexClientSecret)

	resp, err := http.Post(
		"https://oauth.yandex.ru/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}

	json.NewDecoder(resp.Body).Decode(&result)
	return result.AccessToken, nil
}

func (s *Service) GetYandexUserInfo(accessToken string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", "https://login.yandex.ru/info?format=json", nil)
	req.Header.Set("Authorization", "OAuth "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)
	return userInfo, nil
}
