package application

import (
	"encoding/json"
	"errors"
	"financial_data/internal/domain"
	"financial_data/internal/infrastructure"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CompanyHandler struct {
	repo          CompanyRepository
	marketService domain.MarketService
}

func NewCompanyHandler(repo CompanyRepository, marketService domain.MarketService) *CompanyHandler {
	return &CompanyHandler{repo: repo, marketService: marketService}
}

func RegisterCompanyRoutes(r chi.Router, repo CompanyRepository, marketService domain.MarketService) {
	handler := NewCompanyHandler(repo, marketService)

	r.Get("/companies", handler.HandleGetAll)
	r.Get("/companies/{ticker}", handler.HandleGetByTicker)
	r.Get("/companies/sector/{sector_id}", handler.HandleGetBySector)
	r.Post("/companies", handler.HandleCreate)
	r.Put("/companies/{ticker}", handler.HandleUpdate)
	r.Delete("/companies/{ticker}", handler.HandleDelete)
}

func (h *CompanyHandler) HandleGetByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	company, err := h.repo.GetByTicker(r.Context(), ticker)
	if err != nil {
		var dbErr *infrastructure.DbError
		if errors.As(err, &dbErr) && dbErr.RowsAffected == 0 {
			RespondWithError(w, r, 404, dbErr.Message, err)
			return
		}
		RespondWithError(w, r, 500, "failed to load company", err)
		return
	}

	RespondWithSuccess(w, 200, company, "Successfully retrieved company")
}

func (h *CompanyHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	companies, err := h.repo.GetAll(r.Context())
	if err != nil {
		RespondWithError(w, r, 500, "failed to load companies", err)
		return
	}

	RespondWithSuccess(w, 200, companies, "Successfully retrieved companies")
}

func (h *CompanyHandler) HandleGetBySector(w http.ResponseWriter, r *http.Request) {
	sectorIDStr := chi.URLParam(r, "sector_id")
	if sectorIDStr == "" {
		RespondWithError(w, r, 400, "sector_id is required", nil)
		return
	}

	sectorID, err := strconv.Atoi(sectorIDStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid sector_id", err)
		return
	}

	sector := domain.Sector(sectorID)
	if !sector.IsValid() {
		RespondWithError(w, r, 400, "invalid sector_id (allowed values from 1 to 19)", nil)
		return
	}

	companies, err := h.repo.GetBySector(r.Context(), sectorID)
	if err != nil {
		RespondWithError(w, r, 500, "failed to load companies by sector", err)
		return
	}

	RespondWithSuccess(w, 200, companies, "Successfully retrieved companies by sector")
}

func (h *CompanyHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	var company domain.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	sector := domain.Sector(company.SectorID)
	if !sector.IsValid() {
		RespondWithError(w, r, 400, "invalid sector_id (allowed values from 1 to 19)", nil)
		return
	}

	stockInfo, err := h.marketService.GetStockInfo(company.Ticker)
	if err != nil {
		RespondWithError(w, r, 400, "Company with this ticker is not traded on MOEX", err)
		return
	}

	if stockInfo != nil {
		company.Name = stockInfo.Name
	}

	if err := h.repo.Create(r.Context(), &company); err != nil {
		RespondWithError(w, r, 500, "failed to create company", err)
		return
	}

	RespondWithSuccess(w, 201, company, "Company successfully created")
}

func (h *CompanyHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	var company domain.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	sector := domain.Sector(company.SectorID)
	if !sector.IsValid() {
		RespondWithError(w, r, 400, "invalid sector_id (allowed values from 1 to 19)", nil)
		return
	}

	if err := h.repo.Update(r.Context(), ticker, &company); err != nil {
		var dbErr *infrastructure.DbError
		if errors.As(err, &dbErr) && dbErr.RowsAffected == 0 {
			RespondWithError(w, r, 404, dbErr.Message, err)
			return
		}
		RespondWithError(w, r, 500, "failed to update company", err)
		return
	}

	RespondWithSuccess(w, 200, map[string]string{"status": "updated"}, "Company successfully updated")
}

func (h *CompanyHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	if err := h.repo.Delete(r.Context(), ticker); err != nil {
		var dbErr *infrastructure.DbError
		if errors.As(err, &dbErr) && dbErr.RowsAffected == 0 {
			RespondWithError(w, r, 404, dbErr.Message, err)
			return
		}
		RespondWithError(w, r, 500, "failed to delete company", err)
		return
	}

	RespondWithSuccess(w, 204, nil, "Company successfully deleted")
}
