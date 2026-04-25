package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"product-aggregator/internal/model"
	"product-aggregator/internal/provider"
	"product-aggregator/internal/service"
	"product-aggregator/internal/storage"
)

type mockProvider struct {
	name     string
	products []model.Product
	err      error
	delay    time.Duration
	timeout  time.Duration
}

func (m *mockProvider) Name() string           { return m.name }
func (m *mockProvider) Timeout() time.Duration { return m.timeout }
func (m *mockProvider) Fetch(ctx context.Context, _ string) ([]model.Product, error) {
	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.products, nil
}

var _ provider.Provider = (*mockProvider)(nil)

func TestSearchSortsByPrice(t *testing.T) {
	store := storage.New()
	providers := []provider.Provider{
		&mockProvider{
			name:    "p1",
			timeout: 5 * time.Second,
			products: []model.Product{
				{SKU: "A", Name: "Expensive", Price: 50.0, Currency: "USD", Provider: "p1"},
			},
		},
		&mockProvider{
			name:    "p2",
			timeout: 5 * time.Second,
			products: []model.Product{
				{SKU: "B", Name: "Cheap", Price: 10.0, Currency: "USD", Provider: "p2"},
			},
		},
	}
	svc := service.New(providers, store)
	resp := svc.Search(context.Background(), "test")

	if len(resp.Products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(resp.Products))
	}
	if resp.Products[0].Price > resp.Products[1].Price {
		t.Error("products not sorted by price ascending")
	}
}

func TestSearchDeduplicatesBySKU(t *testing.T) {
	store := storage.New()
	providers := []provider.Provider{
		&mockProvider{
			name:    "p1",
			timeout: 5 * time.Second,
			products: []model.Product{
				{SKU: "SAME", Name: "Product 1", Price: 20.0, Currency: "USD", Provider: "p1"},
			},
		},
		&mockProvider{
			name:    "p2",
			timeout: 5 * time.Second,
			products: []model.Product{
				{SKU: "SAME", Name: "Product 2", Price: 15.0, Currency: "USD", Provider: "p2"},
			},
		},
	}
	svc := service.New(providers, store)
	resp := svc.Search(context.Background(), "test")

	if len(resp.Products) != 1 {
		t.Fatalf("expected 1 product after dedup, got %d", len(resp.Products))
	}
}

func TestSearchHandlesProviderFailure(t *testing.T) {
	store := storage.New()
	providers := []provider.Provider{
		&mockProvider{
			name:    "good",
			timeout: 5 * time.Second,
			products: []model.Product{
				{SKU: "A", Name: "Product", Price: 10.0, Currency: "USD", Provider: "good"},
			},
		},
		&mockProvider{
			name:    "bad",
			timeout: 5 * time.Second,
			err:     errors.New("provider error"),
		},
	}
	svc := service.New(providers, store)
	resp := svc.Search(context.Background(), "test")

	if len(resp.Products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(resp.Products))
	}
	if len(resp.FailedProviders) != 1 || resp.FailedProviders[0] != "bad" {
		t.Error("expected 'bad' in failed providers")
	}
}

func TestSearchHandlesTimeout(t *testing.T) {
	store := storage.New()
	providers := []provider.Provider{
		&mockProvider{
			name:    "fast",
			timeout: 5 * time.Second,
			products: []model.Product{
				{SKU: "A", Name: "Product", Price: 10.0, Currency: "USD", Provider: "fast"},
			},
		},
		&mockProvider{
			name:    "slow",
			timeout: 100 * time.Millisecond,
			delay:   2 * time.Second,
			products: []model.Product{
				{SKU: "B", Name: "Slow Product", Price: 5.0, Currency: "USD", Provider: "slow"},
			},
		},
	}
	svc := service.New(providers, store)
	resp := svc.Search(context.Background(), "test")

	if len(resp.Products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(resp.Products))
	}
	if len(resp.FailedProviders) != 1 || resp.FailedProviders[0] != "slow" {
		t.Error("expected 'slow' in failed providers")
	}
}

func TestSearchSavesHistory(t *testing.T) {
	store := storage.New()
	providers := []provider.Provider{
		&mockProvider{
			name:    "p1",
			timeout: 5 * time.Second,
			products: []model.Product{
				{SKU: "A", Name: "Product", Price: 10.0, Currency: "USD", Provider: "p1"},
			},
		},
	}
	svc := service.New(providers, store)
	svc.Search(context.Background(), "mouse")

	history := svc.GetHistory()
	if len(history) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(history))
	}
	if history[0].Query != "mouse" {
		t.Errorf("expected query 'mouse', got '%s'", history[0].Query)
	}
	if history[0].ResultCount != 1 {
		t.Errorf("expected result count 1, got %d", history[0].ResultCount)
	}
}
