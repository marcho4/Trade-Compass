package internal

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router *chi.Mux
}

func CreateServer() *Server {
	log.Println("Creating server")
	return &Server{
		Router: chi.NewRouter(),
	}
}

func (s *Server) MountRoutes() {
	log.Println("Mounting routes")
	s.Router.Get("/webhook", HandleWebhook)
	log.Println("Routes mounted")
}

func (s *Server) Start() {
	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", s.Router)
	if err != nil {
		log.Println("Error starting server: ", err)
	}
	log.Println("Server started")
}
