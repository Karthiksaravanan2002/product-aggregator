package service

import (
	"context"
	"slices"
	"time"

	"product-aggregator/internal/model"
	"product-aggregator/internal/provider"
	"product-aggregator/internal/storage"
)

type Service struct {
	providers []provider.Provider
	storage   *storage.Storage
}

func New(providers []provider.Provider, store *storage.Storage) *Service {
	return &Service{
		providers: providers,
		storage:   store,
	}
}

type providerResult struct {
	products []model.Product
	provider string
	err      error
}

func (s *Service) Search(ctx context.Context, query string) model.SearchResponse {
	results := make(chan providerResult, len(s.providers))

	for _, p := range s.providers {
		go func(p provider.Provider) {
			pCtx, cancel := context.WithTimeout(ctx, p.Timeout())
			defer cancel()
			products, err := p.Fetch(pCtx, query)
			results <- providerResult{products: products, provider: p.Name(), err: err}
		}(p)
	}

	var allProducts []model.Product
	var failedProviders []string

	for range s.providers {
		r := <-results
		if r.err != nil {
			failedProviders = append(failedProviders, r.provider)
			continue
		}
		allProducts = append(allProducts, r.products...)
	}

	allProducts = deduplicate(allProducts)
	slices.SortFunc(allProducts, sortByPriceAsc)

	if allProducts == nil {
		allProducts = []model.Product{}
	}
	if failedProviders == nil {
		failedProviders = []string{}
	}

	resp := model.SearchResponse{
		Query:           query,
		Products:        allProducts,
		TotalProviders:  len(s.providers),
		FailedProviders: failedProviders,
	}

	s.storage.SaveSearch(model.SearchHistory{
		Query:           query,
		Timestamp:       time.Now(),
		ResultCount:     len(allProducts),
		FailedProviders: failedProviders,
	})

	return resp
}

func deduplicate(products []model.Product) []model.Product {
	seen := make(map[string]bool)
	var result []model.Product
	for _, p := range products {
		if !seen[p.SKU] {
			seen[p.SKU] = true
			result = append(result, p)
		}
	}
	return result
}

func sortByPriceAsc(a, b model.Product) int {
	if a.Price < b.Price {
		return -1
	}
	if a.Price > b.Price {
		return 1
	}
	return 0
}

func (s *Service) GetHistory() []model.SearchHistory {
	return s.storage.GetRecentSearches(10)
}

func (s *Service) GetHistoryByID(id int) (model.SearchHistory, bool) {
	return s.storage.GetSearchByID(id)
}
