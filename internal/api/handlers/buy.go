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
		SendErrors(w, apperr.ErrMecrhNameIsEmpty)
		return
	}

	err := h.Stor.BuyItem(r.Context(), userID, item)
	if err != nil {
		SendErrors(w, err)
		h.Logger.Sugar.Infow("internal error", "error: ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
