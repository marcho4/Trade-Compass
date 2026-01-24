package financial_data

import (
	"context"
	"financial_data/internal/application"
	"financial_data/internal/infrastructure"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	ctx := context.Background()
	pool, err := infrastructure.NewPostgresPool(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	ratiosRepo := infrastructure.NewRatiosRepository(pool)
	priceProvider := infrastructure.NewMoexPriceProvider()

	// TODO: Add routes for macro data and news when handlers are implemented
	// macroDataProvider := infrastructure.NewMacroDataProvider(pool)
	// newsProvider := infrastructure.NewNewsProvider()
	// rawDataProvider := infrastructure.NewRawDataProvider()

	r := chi.NewMux()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	application.RegisterRatiosRoutes(r, ratiosRepo)
	application.RegisterPriceRoutes(r, priceProvider)

	log.Println("Starting server on port 8082")
	if err := http.ListenAndServe(":8082", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
