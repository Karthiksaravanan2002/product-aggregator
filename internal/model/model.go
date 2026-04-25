package model

import "time"

type Product struct {
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Provider string  `json:"provider"`
}

type SearchResponse struct {
	Query           string    `json:"query"`
	Products        []Product `json:"products"`
	TotalProviders  int       `json:"-"`
	FailedProviders []string  `json:"failed_providers"`
}

type SearchHistory struct {
	ID              int       `json:"id"`
	Query           string    `json:"query"`
	Timestamp       time.Time `json:"timestamp"`
	ResultCount     int       `json:"result_count"`
	FailedProviders []string  `json:"failed_providers"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}
