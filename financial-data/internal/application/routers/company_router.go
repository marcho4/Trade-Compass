package routers

import (
	"encoding/json"
	"errors"
	"financial_data/internal/application/middleware"
	"financial_data/internal/application/response"
	"financial_data/internal/domain"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CompanyHandler struct {
	repo           CompanyRepository
	marketService  domain.MarketService
	eventPublisher EventPublisher
}

func NewCompanyHandler(repo CompanyRepository, marketService domain.MarketService, eventPublisher EventPublisher) *CompanyHandler {
	return &CompanyHandler{repo: repo, marketService: marketService, eventPublisher: eventPublisher}
}

func RegisterCompanyRoutes(r chi.Router, repo CompanyRepository, marketService domain.MarketService, eventPublisher EventPublisher, m *middleware.MiddlewareConfig) {
	handler := NewCompanyHandler(repo, marketService, eventPublisher)

	r.Get("/companies", handler.HandleGetAll)
	r.Get("/companies/{ticker}", handler.HandleGetByTicker)
	r.Get("/companies/sector/{sector_id}", handler.HandleGetBySector)

	r.Group(func(protected chi.Router) {
		protected.Use(m.AuthMiddleware)

		protected.Post("/companies", handler.HandleCreate)
		protected.Put("/companies/{ticker}", handler.HandleUpdate)
		protected.Delete("/companies/{ticker}", handler.HandleDelete)
	})
}

func (h *CompanyHandler) HandleGetByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	company, err := h.repo.GetByTicker(r.Context(), ticker)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.RespondWithError(w, r, 404, "company not found", err)
			return
		}
		response.RespondWithError(w, r, 500, "failed to load company", err)
		return
	}

	response.RespondWithSuccess(w, 200, company, "Successfully retrieved company")
}

func (h *CompanyHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	companies, err := h.repo.GetAll(r.Context())
	if err != nil {
		response.RespondWithError(w, r, 500, "failed to load companies", err)
		return
	}

	response.RespondWithSuccess(w, 200, companies, "Successfully retrieved companies")
}

func (h *CompanyHandler) HandleGetBySector(w http.ResponseWriter, r *http.Request) {
	sectorIDStr := chi.URLParam(r, "sector_id")
	if sectorIDStr == "" {
		response.RespondWithError(w, r, 400, "sector_id is required", nil)
		return
	}

	sectorID, err := strconv.Atoi(sectorIDStr)
	if err != nil {
		response.RespondWithError(w, r, 400, "invalid sector_id", err)
		return
	}

	sector := domain.Sector(sectorID)
	if !sector.IsValid() {
		response.RespondWithError(w, r, 400, "invalid sector_id (allowed values from 1 to 19)", nil)
		return
	}

	companies, err := h.repo.GetBySector(r.Context(), sectorID)
	if err != nil {
		response.RespondWithError(w, r, 500, "failed to load companies by sector", err)
		return
	}

	response.RespondWithSuccess(w, 200, companies, "Successfully retrieved companies by sector")
}

func (h *CompanyHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var company domain.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		response.RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	sector := domain.Sector(company.SectorID)
	if !sector.IsValid() {
		response.RespondWithError(w, r, 400, "invalid sector_id (allowed values from 1 to 19)", nil)
		return
	}

	stockInfo, err := h.marketService.GetStockInfo(company.Ticker)
	if err != nil {
		response.RespondWithError(w, r, 400, "Company with this ticker is not traded on MOEX", err)
		return
	}

	if stockInfo != nil {
		company.Name = stockInfo.Name
	}

	if err := h.repo.Create(r.Context(), &company); err != nil {
		response.RespondWithError(w, r, 500, "failed to create company", err)
		return
	}

	id := uuid.New().String()
	if err := h.eventPublisher.PublishCompanyCreated(r.Context(), company.Ticker, company.Name, id); err != nil {
		slog.Error("failed to publish company created event", "ticker", company.Ticker, "error", err)
	}

	if err := h.eventPublisher.PublishBusinessResearchTask(r.Context(), company.Ticker, id); err != nil {
		slog.Error("failed to publish business research task", "ticker", company.Ticker, "error", err)
	}

	if err := h.eventPublisher.PublishExpectRiskAndGrowthAnalysis(r.Context(), company.Ticker, id); err != nil {
		slog.Error("failed to publish expect risk and growht task", "ticker", company.Ticker, "error", err)
	}

	response.RespondWithSuccess(w, 201, company, "Company successfully created")
}

func (h *CompanyHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	var company domain.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		response.RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	sector := domain.Sector(company.SectorID)
	if !sector.IsValid() {
		response.RespondWithError(w, r, 400, "invalid sector_id (allowed values from 1 to 19)", nil)
		return
	}

	if err := h.repo.Update(r.Context(), ticker, &company); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.RespondWithError(w, r, 404, "company not found", err)
			return
		}
		response.RespondWithError(w, r, 500, "failed to update company", err)
		return
	}

	response.RespondWithSuccess(w, 200, nil, "Company successfully updated")
}

func (h *CompanyHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	if err := h.repo.Delete(r.Context(), ticker); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.RespondWithError(w, r, 404, "company not found", err)
			return
		}
		response.RespondWithError(w, r, 500, "failed to delete company", err)
		return
	}

	response.RespondWithSuccess(w, 204, nil, "Company successfully deleted")
}
