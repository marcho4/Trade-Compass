package financial_data

import (
	"context"
	"log"
	"net/http"

	"financial_data/application"
	"financial_data/infrastructure"

	"github.com/go-chi/chi"
)

func main() {
	ctx := context.Background()
	pool, err := infrastructure.NewPostgresPool(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	ratiosRepo := infrastructure.NewRatiosRepository(pool)
	ratiosHandler := application.NewRatiosHandler(ratiosRepo)

	r := chi.NewMux()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	r.Get("/ratios/{ticker}", ratiosHandler.HandleGetRatiosByTicker)

	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
