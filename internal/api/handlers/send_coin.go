package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/plasmatrip/avito_merch/internal/apperr"
	"github.com/plasmatrip/avito_merch/internal/model"
)

func (h *Handlers) SendCoin(w http.ResponseWriter, r *http.Request) {
	var sc model.SendCoinRequest

	if err := json.NewDecoder(r.Body).Decode(&sc); err != nil {
		h.Logger.Sugar.Infow("error in request handler", "error: ", err)
		SendErrors(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(model.ValidLogin{}).(*model.Claims).UserdID

	err := h.Stor.SendCoin(r.Context(), userID, sc)
	if err != nil {
		switch err {
		case apperr.ErrInsufficientFunds:
			SendErrors(w, err.Error(), http.StatusBadRequest)
		case apperr.ErrSenderNotFound:
			SendErrors(w, err.Error(), http.StatusBadRequest)
		case apperr.ErrRecipientNotFound:
			SendErrors(w, err.Error(), http.StatusBadRequest)
		default:
			SendErrors(w, err.Error(), http.StatusInternalServerError)
		}

		h.Logger.Sugar.Infow("send coin error", "error: ", err)
		return
	}
}
