package server

import (
	"fmt"
	"net/http"
	"time"

	"product-aggregator/internal/handler"
	"product-aggregator/internal/logger"
)

func NewMux(h *handler.Handler, basePath string) *http.ServeMux {
	mux := http.NewServeMux()

	routes := []struct {
		pattern string
		handler http.HandlerFunc
	}{
		{"GET " + basePath + "/search", h.Search},
		{"GET " + basePath + "/history/{id}", h.GetHistoryByID},
		{"GET " + basePath + "/history", h.GetHistory},
		{"GET " + basePath + "/health", h.Health},
	}

	for _, r := range routes {
		mux.HandleFunc(r.pattern, r.handler)
		if logger.Log != nil {
			logger.Log.Info("registered route", logger.String("pattern", r.pattern))
		}
	}

	return mux
}

func New(h *handler.Handler, port int, timeout time.Duration, basePath string) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           NewMux(h, basePath),
		ReadHeaderTimeout: timeout,
		WriteTimeout:      timeout,
	}
}
