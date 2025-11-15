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

func HandleLogin(w http.ResponseWriter, r *http.Request)          {}
func HandleYandexLogin(w http.ResponseWriter, r *http.Request)    {}
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request)    {}
func HandleYandexCallback(w http.ResponseWriter, r *http.Request) {}
func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {}

func HandlePasswordReset(w http.ResponseWriter, r *http.Request)  {}
func HandleForgotPassword(w http.ResponseWriter, r *http.Request) {}

func HandleLogout(w http.ResponseWriter, r *http.Request) {}
func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {}
func HandleMe(w http.ResponseWriter, r *http.Request) {}