package router

import (
	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
	"github.com/plasmatrip/avito_merch/internal/api"
	"github.com/plasmatrip/avito_merch/internal/api/middleware"
	"github.com/plasmatrip/avito_merch/internal/config"
	"github.com/plasmatrip/avito_merch/internal/logger"
	"github.com/plasmatrip/avito_merch/internal/storage"
)

// NewRouter создает новый маршрутизатор
func NewRouter(cfg config.Config, log logger.Logger, stor storage.Repository, api api.API) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.WithLogging(log), middleware.WithCompression(log), chimid.RedirectSlashes, middleware.WithLimitter(log))

	r.Post("/api/auth", api.Auth)
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.WithAuthentication(log, cfg.TokenSecret))
		r.Get("/info", api.Info)
		r.Post("/sendCoin", api.SendCoin)
		r.Get("/buy/{item}", api.Buy)
	})

	return r
}
