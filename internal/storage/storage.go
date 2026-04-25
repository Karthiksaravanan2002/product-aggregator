package storage

import (
	"sync"

	"product-aggregator/internal/model"
)

type Storage struct {
	mu       sync.RWMutex
	searches []model.SearchHistory
	nextID   int
}

func New() *Storage {
	return &Storage{nextID: 1}
}

func (s *Storage) SaveSearch(h model.SearchHistory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	h.ID = s.nextID
	s.nextID++
	s.searches = append(s.searches, h)
}

func (s *Storage) GetRecentSearches(limit int) []model.SearchHistory {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := len(s.searches)
	if n == 0 {
		return []model.SearchHistory{}
	}

	start := n - limit
	if start < 0 {
		start = 0
	}

	result := make([]model.SearchHistory, n-start)
	copy(result, s.searches[start:])

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func (s *Storage) GetSearchByID(id int) (model.SearchHistory, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, h := range s.searches {
		if h.ID == id {
			return h, true
		}
	}
	return model.SearchHistory{}, false
}
