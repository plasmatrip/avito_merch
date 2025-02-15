package handlers

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/plasmatrip/avito_merch/internal/model"
)

func (h *Handlers) Info(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.ValidLogin{}).(*model.Claims).UserdID

	infoResponse, err := h.Stor.Info(r.Context(), userID)
	if err != nil {
		SendErrors(w, err)
		h.Logger.Sugar.Infow("internal error", "error: ", err)
		return
	}

	err = jsoniter.NewEncoder(w).Encode(infoResponse)
	if err != nil {
		h.Logger.Sugar.Infow("error in request handler", "error: ", err)
		SendErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
