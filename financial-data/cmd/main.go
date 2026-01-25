package financial_data

import (
	"context"
	"financial_data/internal/application"
	"financial_data/internal/infrastructure"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	ctx := context.Background()

	if err := infrastructure.RunMigrations(); err != nil {
		slog.Error("Failed to run migrations", "error", err)
	}

	pool, err := infrastructure.NewPostgresPool(ctx)
	if err != nil {
		slog.Error("Failed to connect to postgres", "error", err)
	}
	defer pool.Close()

	ratiosRepo := infrastructure.NewRatiosRepository(pool)
	rawDataRepo := infrastructure.NewRawDataRepository(pool)
	companyRepo := infrastructure.NewCompanyRepository(pool)
	sectorRepo := infrastructure.NewSectorRepository(pool)
	dividendsRepo := infrastructure.NewDividendsRepository(pool)
	cbRateRepo := infrastructure.NewCBRateRepository(pool)
	newsRepo := infrastructure.NewNewsRepository(pool)

	priceProvider := infrastructure.NewMoexPriceProvider()

	r := chi.NewMux()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	application.RegisterRatiosRoutes(r, ratiosRepo)
	application.RegisterRawDataRoutes(r, rawDataRepo)
	application.RegisterCompanyRoutes(r, companyRepo)
	application.RegisterSectorRoutes(r, sectorRepo)
	application.RegisterDividendsRoutes(r, dividendsRepo)
	application.RegisterMacroRoutes(r, cbRateRepo)
	application.RegisterNewsRoutes(r, newsRepo)
	application.RegisterPriceRoutes(r, priceProvider)

	if err := http.ListenAndServe(":8082", r); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
