package handlers

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/plasmatrip/avito_merch/internal/apperr"
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

func SendErrors(w http.ResponseWriter, err error) {
	msg, ok := apperr.ErrorMessages[err]
	if !ok {
		msg = "internal error"
	}
	statusCode, ok := apperr.ErrorStatuses[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errResponse := model.ErrorResponse{Errors: msg}
	jsoniter.NewEncoder(w).Encode(errResponse)
}
