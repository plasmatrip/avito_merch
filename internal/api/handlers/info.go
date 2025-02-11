package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/plasmatrip/avito_merch/internal/model"
)

func (h *Handlers) Info(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.ValidLogin{}).(*model.Claims).UserdID

	infoResponse, err := h.Stor.Info(r.Context(), userID)
	if err != nil {
		h.Logger.Sugar.Infow("info error", "error: ", err)
		SendErrors(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(infoResponse)
}
