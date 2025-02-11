package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/plasmatrip/avito_merch/internal/api/handlers"
	"github.com/plasmatrip/avito_merch/internal/api/middleware"
	"github.com/plasmatrip/avito_merch/internal/config"
	"github.com/plasmatrip/avito_merch/internal/logger"
	"github.com/plasmatrip/avito_merch/internal/storage"
)

// NewRouter создает новый маршрутизатор
func NewRouter(cfg config.Config, log logger.Logger, stor storage.Repository) *chi.Mux {

	r := chi.NewRouter()

	handlers := handlers.Handlers{Config: cfg, Logger: log, Stor: stor}

	r.Use(middleware.WithLogging(log), middleware.WithCompression(log))

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/", handlers.Auth)
	})

	r.Route("/api/info", func(r chi.Router) {
		r.Use(middleware.WithAuthentication(log, cfg.TokenSecret))
		r.Get("/", handlers.Info)
	})

	r.Route("/api/sendCoin", func(r chi.Router) {
		r.Use(middleware.WithAuthentication(log, cfg.TokenSecret))
		r.Post("/", handlers.SendCoin)
	})

	r.Route("/api/buy/{item}", func(r chi.Router) {
		r.Use(middleware.WithAuthentication(log, cfg.TokenSecret))
		r.Get("/", handlers.Buy)
	})

	return r
}
