package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"product-aggregator/internal/handler"
	"product-aggregator/internal/model"
	"product-aggregator/internal/provider"
	"product-aggregator/internal/server"
	"product-aggregator/internal/service"
	"product-aggregator/internal/storage"
)

type mockProvider struct {
	name     string
	products []model.Product
}

func (m *mockProvider) Name() string           { return m.name }
func (m *mockProvider) Timeout() time.Duration { return 5 * time.Second }
func (m *mockProvider) Fetch(_ context.Context, _ string) ([]model.Product, error) {
	return m.products, nil
}

var _ provider.Provider = (*mockProvider)(nil)

func setupMux() *http.ServeMux {
	store := storage.New()
	providers := []provider.Provider{
		&mockProvider{
			name: "test_a",
			products: []model.Product{
				{SKU: "SKU1", Name: "Mouse", Price: 25.99, Currency: "USD", Provider: "test_a"},
			},
		},
		&mockProvider{
			name: "test_b",
			products: []model.Product{
				{SKU: "SKU2", Name: "Keyboard", Price: 15.99, Currency: "USD", Provider: "test_b"},
			},
		},
	}
	svc := service.New(providers, store)
	h := handler.New(svc)
	return server.NewMux(h, "")
}

func TestSearchHandler(t *testing.T) {
	mux := setupMux()
	req := httptest.NewRequest("GET", "/search?q=test", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp model.SearchResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Query != "test" {
		t.Errorf("expected query 'test', got '%s'", resp.Query)
	}
}

func TestSearchHandlerMissingQuery(t *testing.T) {
	mux := setupMux()
	req := httptest.NewRequest("GET", "/search", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHealthHandler(t *testing.T) {
	mux := setupMux()
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp model.HealthResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", resp.Status)
	}
}

func TestHistoryHandler(t *testing.T) {
	mux := setupMux()

	req := httptest.NewRequest("GET", "/search?q=test", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	req = httptest.NewRequest("GET", "/history", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var history []model.SearchHistory
	json.NewDecoder(w.Body).Decode(&history)
	if len(history) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(history))
	}
}

func TestHistoryByIDHandler(t *testing.T) {
	mux := setupMux()

	req := httptest.NewRequest("GET", "/search?q=test", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	req = httptest.NewRequest("GET", "/history/1", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var h model.SearchHistory
	json.NewDecoder(w.Body).Decode(&h)
	if h.Query != "test" {
		t.Errorf("expected query 'test', got '%s'", h.Query)
	}
}

func TestHistoryByIDNotFound(t *testing.T) {
	mux := setupMux()
	req := httptest.NewRequest("GET", "/history/999", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
