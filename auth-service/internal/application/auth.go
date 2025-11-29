package internal

import (
	"log"
	"net/http"
)

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	log.Println("Got webhook request")
	w.Write([]byte("Hello, World!!"))
}

func HandleRegistration(w http.ResponseWriter, r *http.Request)      {}
func HandleEmailVerification(w http.ResponseWriter, r *http.Request) {}
func HandleEmailResend(w http.ResponseWriter, r *http.Request)       {}

// func HandleLogin(w http.ResponseWriter, r *http.Request) {
// 	// json.Unmarshal(r.Body, &LoginRequest{})
// 	body, err := io.ReadAll(r.Body)
// 	defer r.Body.Close()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.Unmarshal(body, &LoginRequest{})
// 	db.GetAuthDataByEmail(loginRequest.Email)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if authData.PasswordHash != loginRequest.Password {
// 		http.Error(w, "Invalid password", http.StatusUnauthorized)
// 		return
// 	}

// 	accessToken, refreshToken, err := jwt.GenerateTokens(authData.UserID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"accessToken":  accessToken,
// 		"refreshToken": refreshToken,
// 	})
// 	return
// }

func HandleYandexLogin(w http.ResponseWriter, r *http.Request) {}
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {}

func HandleYandexCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("Got yandex callback request")
	log.Println(r.URL.Query())
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {}

func HandlePasswordReset(w http.ResponseWriter, r *http.Request)  {}
func HandleForgotPassword(w http.ResponseWriter, r *http.Request) {}

func HandleLogout(w http.ResponseWriter, r *http.Request)       {}
func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {}
func HandleMe(w http.ResponseWriter, r *http.Request)           {}
