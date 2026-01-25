package application

import (
	"encoding/json"
	"financial_data/internal/domain"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type CompanyHandler struct {
	repo CompanyRepository
}

func NewCompanyHandler(repo CompanyRepository) *CompanyHandler {
	return &CompanyHandler{repo: repo}
}

func RegisterCompanyRoutes(r chi.Router, repo CompanyRepository) {
	handler := NewCompanyHandler(repo)

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
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	company, err := h.repo.GetByTicker(r.Context(), ticker)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load company: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(company); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *CompanyHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	companies, err := h.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load companies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(companies); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *CompanyHandler) HandleGetBySector(w http.ResponseWriter, r *http.Request) {
	sectorIDStr := chi.URLParam(r, "sector_id")
	if sectorIDStr == "" {
		http.Error(w, "sector_id is required", http.StatusBadRequest)
		return
	}

	sectorID, err := strconv.Atoi(sectorIDStr)
	if err != nil {
		http.Error(w, "invalid sector_id", http.StatusBadRequest)
		return
	}

	sector := domain.Sector(sectorID)
	if !sector.IsValid() {
		http.Error(w, "invalid sector_id (allowed values from 1 to 19)", http.StatusBadRequest)
		return
	}

	companies, err := h.repo.GetBySector(r.Context(), sectorID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load companies by sector: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(companies); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *CompanyHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var company domain.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sector := domain.Sector(company.SectorID)
	if !sector.IsValid() {
		http.Error(w, "invalid sector_id (allowed values from 1 to 19)", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &company); err != nil {
		http.Error(w, fmt.Sprintf("failed to create company: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

func (h *CompanyHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	var company domain.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sector := domain.Sector(company.SectorID)
	if !sector.IsValid() {
		http.Error(w, "invalid sector_id (allowed values from 1 to 19)", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(r.Context(), ticker, &company); err != nil {
		http.Error(w, fmt.Sprintf("failed to update company: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *CompanyHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), ticker); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete company: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
