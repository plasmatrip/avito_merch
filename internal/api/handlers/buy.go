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

		switch err {
		case apperr.ErrItemNotFound:
			h.Logger.Sugar.Infow("buy error", "error: ", err)
			SendErrors(w, "item not found", http.StatusBadRequest)
		case apperr.ErrInsufficientFunds:
			h.Logger.Sugar.Infow("buy error", "error: ", err)
			SendErrors(w, "insufficient funds", http.StatusBadRequest)
		case apperr.ErrAccountNotFound:
			h.Logger.Sugar.Infow("buy error", "error: ", err)
			SendErrors(w, "account not found", http.StatusBadRequest)
		case apperr.ErrMerchNotBought:
			h.Logger.Sugar.Infow("buy error", "error: ", err)
			SendErrors(w, "merch not bought", http.StatusInternalServerError)
		default:
			h.Logger.Sugar.Infow("buy error", "error: ", err)
			SendErrors(w, "buy error", http.StatusInternalServerError)
		}

		return
	}

	h.Logger.Sugar.Infow("buy request", "userID: ", userID, "body: ", item)
}
