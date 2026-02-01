package application

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Status  string `json:"status"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func RespondWithError(w http.ResponseWriter, r *http.Request, code int, message string, err error) {
	if err != nil {
		slog.Error("Request failed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", code,
			"error", err,
		)
	}

	RespondWithJSON(w, code, ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Failed to marshal JSON response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithSuccess(w http.ResponseWriter, code int, data interface{}, message string) {
	response := SuccessResponse{
		Status:  "success",
		Data:    data,
		Message: message,
	}
	RespondWithJSON(w, code, response)
}
