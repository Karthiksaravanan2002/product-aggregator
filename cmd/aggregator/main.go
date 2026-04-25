package main

import (
	"time"

	"product-aggregator/internal/handler"
	"product-aggregator/internal/logger"
	"product-aggregator/internal/props"
	"product-aggregator/internal/provider"
	"product-aggregator/internal/server"
	"product-aggregator/internal/service"
	"product-aggregator/internal/storage"
)

func main() {
	logger.Init()
	defer logger.Sync()

	cfg, err := props.Load("env/application.yaml")
	if err != nil {
		logger.Log.Fatal("failed to load config", logger.Error(err))
	}

	store := storage.New()

	registry := provider.NewRegistry()
	registry.Register(cfg.Provider.A.Enabled, provider.NewProviderA(time.Duration(cfg.Provider.A.TimeoutMs)*time.Millisecond))
	registry.Register(cfg.Provider.B.Enabled, provider.NewProviderB(time.Duration(cfg.Provider.B.TimeoutMs)*time.Millisecond))
	registry.Register(cfg.Provider.C.Enabled, provider.NewProviderC(time.Duration(cfg.Provider.C.TimeoutMs)*time.Millisecond))

	svc := service.New(registry.All(), store)
	h := handler.New(svc)
	srvTimeout := time.Duration(cfg.Server.TimeoutMs) * time.Millisecond
	srv := server.New(h, cfg.Server.Port, srvTimeout, cfg.Server.BasePath)

	logger.Log.Info("starting server", logger.Int("port", cfg.Server.Port))
	if err := srv.ListenAndServe(); err != nil {
		logger.Log.Fatal("server failed", logger.Error(err))
	}
}
