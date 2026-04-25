package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"product-aggregator/internal/model"
	"product-aggregator/internal/service"
)

type Handler struct {
	service *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "query parameter 'q' is required"})
		return
	}
	result := h.service.Search(r.Context(), query)

	status := http.StatusOK
	failed := len(result.FailedProviders)
	if failed == result.TotalProviders {
		status = http.StatusBadGateway
	} else if failed > 0 {
		status = http.StatusMultiStatus
	}

	writeJSON(w, status, result)
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	history := h.service.GetHistory()
	writeJSON(w, http.StatusOK, history)
}

func (h *Handler) GetHistoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	history, found := h.service.GetHistoryByID(id)
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "search not found"})
		return
	}
	writeJSON(w, http.StatusOK, history)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, model.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
