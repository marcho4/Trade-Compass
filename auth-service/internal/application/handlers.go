package internal

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"auth-service/internal/config"
	"auth-service/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type Handlers struct {
	Service *Service
	Config  *config.Config
}

func NewHandlers(service *Service, cfg *config.Config) *Handlers {
	return &Handlers{
		Service: service,
		Config:  cfg,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Email  string `json:"email,omitempty"`
}

func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handlers) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		respondError(w, http.StatusBadRequest, "Email, password and name are required")
		return
	}

	if len(req.Password) < 8 {
		respondError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to hash password")
		log.Printf("Failed to hash password: %v", err)
		return
	}

	deviceInfo := getDeviceInfo(r)

	ctx := r.Context()
	accessToken, refreshToken, err := h.Service.CreateUser(ctx, req.Name, req.Email, hashedPassword, deviceInfo)
	if err != nil {
		if err.Error() == "user already exists" {
			respondError(w, http.StatusConflict, "User already exists")
			return
		}
		log.Printf("Error creating user: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	h.setAuthCookies(w, r, accessToken, refreshToken)
	respondJSON(w, http.StatusCreated, map[string]string{"message": "User created"})
}

func (h *Handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	deviceInfo := getDeviceInfo(r)

	ctx := r.Context()
	accessToken, refreshToken, err := h.Service.ProcessLogin(ctx, req.Email, req.Password, deviceInfo)
	if err != nil {
		if err.Error() == "invalid credentials" {
			respondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		log.Printf("Error processing login: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to process login")
		return
	}

	h.setAuthCookies(w, r, accessToken, refreshToken)

	respondJSON(w, http.StatusOK, map[string]string{"message": "Login successful"})
}

func (h *Handlers) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	refreshCookie, err := r.Cookie("refreshToken")
	if err != nil || refreshCookie.Value == "" {
		respondError(w, http.StatusUnauthorized, "Refresh token required")
		return
	}

	deviceInfo := getDeviceInfo(r)
	ctx := r.Context()

	accessToken, refreshToken, err := h.Service.RefreshTokens(ctx, refreshCookie.Value, deviceInfo)
	if err != nil {
		log.Printf("Error refreshing tokens: %v", err)
		respondError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	h.setAuthCookies(w, r, accessToken, refreshToken)

	respondJSON(w, http.StatusOK, map[string]string{"message": "Token refreshed"})
}

func (h *Handlers) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	refreshCookie, err := r.Cookie("refreshToken")
	if err == nil && refreshCookie.Value != "" {
		ctx := r.Context()
		if err := h.Service.RevokeRefreshToken(ctx, refreshCookie.Value); err != nil {
			log.Printf("Error revoking token: %v", err)
		}
	}

	h.clearAuthCookies(w, r)

	respondJSON(w, http.StatusOK, map[string]string{"message": "Logged out"})
}

func (h *Handlers) HandleMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	tokenCookie, err := r.Cookie("accessToken")
	if err != nil || tokenCookie.Value == "" {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	claims, err := h.Service.ValidateAccessToken(tokenCookie.Value)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	ctx := r.Context()
	user, err := h.Service.GetUserByID(ctx, claims.UserID)
	if err != nil || user == nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	respondJSON(w, http.StatusOK, UserResponse{
		ID:     user.GetID(),
		Name:   user.GetName(),
		Status: user.GetStatus(),
	})
}

func (h *Handlers) HandleYandexLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateSecureState()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate state")
		return
	}

	isSecure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})

	redirectURL := h.getRedirectURL(r)
	callbackURL := redirectURL + "/api/auth/callback/yandex"

	yandexAuthURL := fmt.Sprintf(
		"https://oauth.yandex.ru/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s",
		h.Config.OAuth.Yandex.ClientID,
		callbackURL,
		state,
	)

	respondJSON(w, http.StatusOK, map[string]string{
		"url": yandexAuthURL,
	})
}

func (h *Handlers) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusNotImplemented, "Google OAuth not implemented yet")
}

func (h *Handlers) HandleYandexCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	redirectURL := h.getRedirectURL(r)

	savedState, err := r.Cookie("oauth_state")
	if err != nil || savedState.Value == "" || savedState.Value != state {
		log.Printf("Yandex callback: invalid or missing state (saved: %v, received: %s)", savedState, state)
		http.Redirect(w, r, redirectURL+"/auth?error=invalid_state", http.StatusFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	if code == "" {
		log.Printf("Yandex callback: code parameter is missing")
		http.Redirect(w, r, redirectURL+"/auth?error=missing_code", http.StatusFound)
		return
	}

	ctx := r.Context()
	deviceInfo := getDeviceInfo(r)

	accessToken, err := h.Service.GetYandexAccessToken(code)
	if err != nil {
		log.Printf("Error getting Yandex access token: %v", err)
		http.Redirect(w, r, redirectURL+"/auth?error=token_failed", http.StatusFound)
		return
	}

	userInfo, err := h.Service.GetYandexUserInfo(accessToken)
	if err != nil {
		log.Printf("Error getting Yandex user info: %v", err)
		http.Redirect(w, r, redirectURL+"/auth?error=user_info_failed", http.StatusFound)
		return
	}

	yandexID, _ := userInfo["id"].(string)
	email, _ := userInfo["default_email"].(string)
	name, _ := userInfo["display_name"].(string)
	if name == "" {
		name, _ = userInfo["real_name"].(string)
	}
	if name == "" {
		name, _ = userInfo["login"].(string)
	}
	if name == "" {
		name = "Yandex User"
	}

	if yandexID == "" || email == "" {
		log.Printf("Invalid user data from Yandex: yandexID=%s, email=%s", yandexID, email)
		http.Redirect(w, r, redirectURL+"/auth?error=invalid_user_data", http.StatusFound)
		return
	}

	providerAuth, err := h.Service.GetProviderAuth(ctx, domain.ProviderTypeYandex, yandexID)
	if err != nil {
		log.Printf("Error getting provider auth: %v", err)
		http.Redirect(w, r, redirectURL+"/auth?error=auth_check_failed", http.StatusFound)
		return
	}

	var user *domain.User
	var accessTokenStr, refreshTokenStr string

	if providerAuth != nil {
		user, err = h.Service.GetUserByID(ctx, providerAuth.GetUserID())
		if err != nil || user == nil {
			log.Printf("Error getting user: %v", err)
			http.Redirect(w, r, redirectURL+"/auth?error=user_not_found", http.StatusFound)
			return
		}
	} else {
		existingUser, err := h.Service.GetUserByEmail(ctx, email)
		if err != nil {
			log.Printf("Error checking existing user: %v", err)
			http.Redirect(w, r, redirectURL+"/auth?error=user_check_failed", http.StatusFound)
			return
		}

		if existingUser != nil {
			user = existingUser
		} else {
			user, err = h.Service.CreateOAuthUser(ctx, name, email)
			if err != nil {
				log.Printf("Error creating user: %v", err)
				http.Redirect(w, r, redirectURL+"/auth?error=user_creation_failed", http.StatusFound)
				return
			}
		}

		if err := h.Service.CreateProviderAuth(ctx, user.GetID(), domain.ProviderTypeYandex, yandexID); err != nil {
			log.Printf("Error creating provider auth: %v", err)
		}
	}

	accessTokenStr, refreshTokenStr, err = h.Service.GenerateTokens(ctx, user.GetID(), user.GetName(), user.GetStatus(), deviceInfo)
	if err != nil {
		log.Printf("Error generating tokens: %v", err)
		http.Redirect(w, r, redirectURL+"/auth?error=token_generation_failed", http.StatusFound)
		return
	}

	if err := h.Service.UpdateLastLogin(ctx, user.GetID()); err != nil {
		log.Printf("Error updating last login: %v", err)
	}

	h.setAuthCookies(w, r, accessTokenStr, refreshTokenStr)

	http.Redirect(w, r, redirectURL+"/welcome", http.StatusFound)
}

func (h *Handlers) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusNotImplemented, "Google OAuth callback not implemented yet")
}

func (h *Handlers) HandlePasswordReset(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusNotImplemented, "Password reset not implemented yet")
}

func (h *Handlers) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusNotImplemented, "Forgot password not implemented yet")
}

func (h *Handlers) HandleEmailVerification(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusNotImplemented, "Email verification not implemented yet")
}

func (h *Handlers) HandleEmailResend(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusNotImplemented, "Email resend not implemented yet")
}

func decodeJSON(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(v)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}

func getDeviceInfo(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "unknown"
	}
	return userAgent
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hashedBytes), nil
}

func generateSecureState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (h *Handlers) setAuthCookies(w http.ResponseWriter, r *http.Request, accessToken, refreshToken string) {
	domain := h.getCookieDomain(r)
	isSecure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"

	accessTokenMaxAge := int(h.Config.JWT.AccessTokenTTL.Seconds())
	refreshTokenMaxAge := int(h.Config.JWT.RefreshTokenTTL.Seconds())

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		Domain:   domain,
		MaxAge:   accessTokenMaxAge,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		Domain:   domain,
		MaxAge:   refreshTokenMaxAge,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handlers) clearAuthCookies(w http.ResponseWriter, r *http.Request) {
	domain := h.getCookieDomain(r)
	isSecure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    "",
		Path:     "/",
		Domain:   domain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/",
		Domain:   domain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handlers) getCookieDomain(r *http.Request) string {
	if h.Config.Frontend.CookieDomain != "" {
		return h.Config.Frontend.CookieDomain
	}

	host := r.Host
	if host == "" {
		host = r.Header.Get("Host")
	}

	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	return host
}

func (h *Handlers) getRedirectURL(r *http.Request) string {
	if h.Config.Frontend.URL != "" {
		return h.Config.Frontend.URL
	}

	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	host := r.Host
	if host == "" {
		host = r.Header.Get("Host")
	}

	return scheme + "://" + host
}
