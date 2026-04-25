package storage_test

import (
	"sync"
	"testing"
	"time"

	"product-aggregator/internal/model"
	"product-aggregator/internal/storage"
)

func TestSaveAndRetrieve(t *testing.T) {
	s := storage.New()
	s.SaveSearch(model.SearchHistory{
		Query:       "test",
		Timestamp:   time.Now(),
		ResultCount: 5,
	})

	results := s.GetRecentSearches(10)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].Query != "test" {
		t.Errorf("expected 'test', got '%s'", results[0].Query)
	}
	if results[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", results[0].ID)
	}
}

func TestGetRecentSearchesLimit(t *testing.T) {
	s := storage.New()
	for i := 0; i < 15; i++ {
		s.SaveSearch(model.SearchHistory{
			Query:       "test",
			Timestamp:   time.Now(),
			ResultCount: i,
		})
	}

	results := s.GetRecentSearches(10)
	if len(results) != 10 {
		t.Fatalf("expected 10, got %d", len(results))
	}
	if results[0].ID != 15 {
		t.Errorf("expected most recent ID 15, got %d", results[0].ID)
	}
}

func TestGetSearchByID(t *testing.T) {
	s := storage.New()
	s.SaveSearch(model.SearchHistory{Query: "first", Timestamp: time.Now()})
	s.SaveSearch(model.SearchHistory{Query: "second", Timestamp: time.Now()})

	h, found := s.GetSearchByID(2)
	if !found {
		t.Fatal("expected to find search with ID 2")
	}
	if h.Query != "second" {
		t.Errorf("expected 'second', got '%s'", h.Query)
	}

	_, found = s.GetSearchByID(99)
	if found {
		t.Error("expected not found for ID 99")
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := storage.New()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.SaveSearch(model.SearchHistory{
				Query:       "concurrent",
				Timestamp:   time.Now(),
				ResultCount: 1,
			})
		}()
	}

	wg.Wait()
	results := s.GetRecentSearches(200)
	if len(results) != 100 {
		t.Fatalf("expected 100, got %d", len(results))
	}
}
