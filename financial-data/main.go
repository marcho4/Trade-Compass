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
	priceProvider := infrastructure.NewMoexPriceProvider()

	ratiosHandler := application.NewRatiosHandler(ratiosRepo)
	priceHandler := application.NewPriceHandler(priceProvider)

	r := chi.NewMux()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/ratios/{ticker}", ratiosHandler.HandleGetRatiosByTicker)
	r.Get("/price", priceHandler.HandleGetPriceByTicker)

	log.Println("Starting server on port 8082")
	if err := http.ListenAndServe(":8082", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
