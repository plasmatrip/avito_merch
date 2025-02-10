package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/plasmatrip/avito_merch/internal/config"
	"github.com/plasmatrip/avito_merch/internal/logger"
	"github.com/plasmatrip/avito_merch/internal/model"
	"github.com/plasmatrip/avito_merch/internal/storage"
)

type Handlers struct {
	Config config.Config
	Logger logger.Logger
	Stor   storage.Repository
}

func NewHandlers(cfg config.Config, l logger.Logger, db storage.Repository) *Handlers {
	return &Handlers{
		Config: cfg,
		Logger: l,
		Stor:   db,
	}
}

func SendErrors(w http.ResponseWriter, error string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errResponse := model.ErrorResponse{Errors: error}
	json.NewEncoder(w).Encode(errResponse)
}
