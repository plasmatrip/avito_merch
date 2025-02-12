package handlers

import (
	"net/http"

	"github.com/plasmatrip/avito_merch/internal/apperr"
	"github.com/plasmatrip/avito_merch/internal/model"
)

func (h *Handlers) Buy(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.ValidLogin{}).(*model.Claims).UserdID

	item := r.PathValue("item")
	if len(item) == 0 {
		h.Logger.Sugar.Infoln("Merch name is empty")
		SendErrors(w, "Merch name is empty", http.StatusBadRequest)
		return
	}

	err := h.Stor.BuyItem(r.Context(), userID, item)
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
	w.WriteHeader(http.StatusOK)
}
