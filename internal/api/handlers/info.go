package handlers

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/plasmatrip/avito_merch/internal/apperr"
	"github.com/plasmatrip/avito_merch/internal/model"
)

func (h *Handlers) Info(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.ValidLogin{}).(*model.Claims).UserdID

	infoResponse, err := h.Stor.Info(r.Context(), userID)
	if err != nil {
		msg, ok := apperr.ErrorMessages[err]
		if !ok {
			msg = "internal error"
		}
		status, ok := apperr.ErrorStatuses[err]
		if !ok {
			status = http.StatusInternalServerError
		}

		SendErrors(w, msg, status)
		h.Logger.Sugar.Infow("internal error", "error: ", err)
		return
	}

	err = jsoniter.NewEncoder(w).Encode(infoResponse)
	if err != nil {
		h.Logger.Sugar.Infow("error in request handler", "error: ", err)
		SendErrors(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
