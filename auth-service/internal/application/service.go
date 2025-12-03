package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

type Service struct {
	dbRepo      domain.DbRepo
	jwtService  domain.JWTService
	oauthConfig config.OAuthConfig
}

func NewService(dbRepo domain.DbRepo, jwtService domain.JWTService, oauthConfig config.OAuthConfig) *Service {
	return &Service{
		dbRepo:      dbRepo,
		jwtService:  jwtService,
		oauthConfig: oauthConfig,
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

func (s *Service) ProcessLogin(ctx context.Context, email, password, deviceInfo string) (string, string, error) {
	authData, err := s.dbRepo.GetAuthDataByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("get auth data: %w", err)
	}

	if authData == nil {
		return "", "", errors.New("invalid credentials")
	}

	if !checkPassword(authData.GetPasswordHash(), password) {
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

func (s *Service) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	return s.dbRepo.GetUserByID(ctx, uint64(userID))
}

func (s *Service) GetProviderAuth(ctx context.Context, providerType domain.ProviderType, providerUserID string) (*domain.ProviderAuth, error) {
	return s.dbRepo.GetProviderAuth(ctx, providerType, providerUserID)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.dbRepo.GetUserByEmail(ctx, email)
}

func (s *Service) CreateOAuthUser(ctx context.Context, name, email string) (*domain.User, error) {
	return s.dbRepo.CreateUser(ctx, name, email, "$oauth$no_password_login_disabled$")
}

func (s *Service) CreateProviderAuth(ctx context.Context, userID int64, providerType domain.ProviderType, providerUserID string) error {
	return s.dbRepo.CreateProviderAuth(ctx, userID, providerType, providerUserID)
}

func (s *Service) UpdateLastLogin(ctx context.Context, userID int64) error {
	return s.dbRepo.UpdateLastLogin(ctx, userID)
}

func (s *Service) GetYandexAccessToken(yandexCode string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", yandexCode)
	data.Set("client_id", s.oauthConfig.Yandex.ClientID)
	data.Set("client_secret", s.oauthConfig.Yandex.ClientSecret)

	resp, err := httpClient.Post(
		"https://oauth.yandex.ru/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("request yandex token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("yandex token request failed with status: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode yandex token response: %w", err)
	}

	if result.Error != "" {
		return "", fmt.Errorf("yandex oauth error: %s", result.Error)
	}

	return result.AccessToken, nil
}

func (s *Service) GetYandexUserInfo(accessToken string) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, "https://login.yandex.ru/info?format=json", nil)
	if err != nil {
		return nil, fmt.Errorf("create yandex user info request: %w", err)
	}
	req.Header.Set("Authorization", "OAuth "+accessToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request yandex user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yandex user info request failed with status: %d", resp.StatusCode)
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("decode yandex user info: %w", err)
	}

	return userInfo, nil
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
