package provider

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"product-aggregator/internal/model"
)

type Provider interface {
	Name() string
	Timeout() time.Duration
	Fetch(ctx context.Context, query string) ([]model.Product, error)
}

type Registry struct {
	providers []Provider
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Register(enabled bool, p Provider) {
	if !enabled {
		return
	}
	r.providers = append(r.providers, p)
}

func (r *Registry) All() []Provider {
	return r.providers
}

func simulateLatency(ctx context.Context) error {
	delay := time.Duration(100+rand.Intn(1100)) * time.Millisecond
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func shouldFail() bool {
	return rand.Intn(5) == 0
}

type ProviderA struct {
	name    string
	timeout time.Duration
}

func NewProviderA(timeout time.Duration) *ProviderA {
	return &ProviderA{name: "provider_a", timeout: timeout}
}

type providerAItem struct {
	SKU         string
	ProductName string
	Price       float64
	Currency    string
}

func (p *ProviderA) Name() string           { return p.name }
func (p *ProviderA) Timeout() time.Duration { return p.timeout }

func (p *ProviderA) Fetch(ctx context.Context, query string) ([]model.Product, error) {
	if err := simulateLatency(ctx); err != nil {
		return nil, err
	}
	if shouldFail() {
		return nil, errors.New("provider_a: internal error")
	}

	catalog := []providerAItem{
		{"SKU-A1", "Wireless Mouse", 25.99, "USD"},
		{"SKU-A2", "Mechanical Keyboard", 79.99, "USD"},
		{"SKU-A3", "USB-C Hub", 34.50, "USD"},
		{"SKU-COMMON1", "Gaming Mouse Pad", 15.99, "USD"},
	}

	q := strings.ToLower(query)
	var products []model.Product
	for _, item := range catalog {
		if strings.Contains(strings.ToLower(item.ProductName), q) {
			products = append(products, model.Product{
				SKU:      item.SKU,
				Name:     item.ProductName,
				Price:    item.Price,
				Currency: item.Currency,
				Provider: p.Name(),
			})
		}
	}
	return products, nil
}

type ProviderB struct {
	name    string
	timeout time.Duration
}

func NewProviderB(timeout time.Duration) *ProviderB {
	return &ProviderB{name: "provider_b", timeout: timeout}
}

type providerBProduct struct {
	ID            string
	Title         string
	Cost          float64
	MoneyCurrency string
}

func (p *ProviderB) Name() string           { return p.name }
func (p *ProviderB) Timeout() time.Duration { return p.timeout }

func (p *ProviderB) Fetch(ctx context.Context, query string) ([]model.Product, error) {
	if err := simulateLatency(ctx); err != nil {
		return nil, err
	}
	if shouldFail() {
		return nil, errors.New("provider_b: Returned error")
	}

	catalog := []providerBProduct{
		{"SKU-B1", "Gaming Mouse", 45.00, "USD"},
		{"SKU-B2", "Laptop Stand", 55.00, "USD"},
		{"SKU-B3", "Wireless Keyboard", 65.00, "USD"},
		{"SKU-COMMON1", "Gaming Mouse Pad", 14.99, "USD"},
	}

	q := strings.ToLower(query)
	var products []model.Product
	for _, item := range catalog {
		if strings.Contains(strings.ToLower(item.Title), q) {
			products = append(products, model.Product{
				SKU:      item.ID,
				Name:     item.Title,
				Price:    item.Cost,
				Currency: item.MoneyCurrency,
				Provider: p.Name(),
			})
		}
	}
	return products, nil
}

type ProviderC struct {
	name    string
	timeout time.Duration
}

func NewProviderC(timeout time.Duration) *ProviderC {
	return &ProviderC{name: "provider_c", timeout: timeout}
}

type providerCResult struct {
	ProductSKU   string
	ProductName  string
	Amount       float64
	CurrencyCode string
}

func (p *ProviderC) Name() string           { return p.name }
func (p *ProviderC) Timeout() time.Duration { return p.timeout }

func (p *ProviderC) Fetch(ctx context.Context, query string) ([]model.Product, error) {
	if err := simulateLatency(ctx); err != nil {
		return nil, err
	}
	if shouldFail() {
		return nil, errors.New("provider_c: timeout")
	}

	catalog := []providerCResult{
		{"SKU-C1", "Bluetooth Mouse", 19.99, "USD"},
		{"SKU-C2", "Monitor Arm", 89.99, "USD"},
		{"SKU-C3", "Desk Pad", 22.50, "USD"},
		{"SKU-COMMON1", "Gaming Mouse Pad", 16.49, "USD"},
	}

	q := strings.ToLower(query)
	var products []model.Product
	for _, item := range catalog {
		if strings.Contains(strings.ToLower(item.ProductName), q) {
			products = append(products, model.Product{
				SKU:      item.ProductSKU,
				Name:     item.ProductName,
				Price:    item.Amount,
				Currency: item.CurrencyCode,
				Provider: p.Name(),
			})
		}
	}
	return products, nil
}
