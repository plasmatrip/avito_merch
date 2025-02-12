package handlers

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/plasmatrip/avito_merch/internal/apperr"
	"github.com/plasmatrip/avito_merch/internal/model"
)

func (h *Handlers) SendCoin(w http.ResponseWriter, r *http.Request) {
	var sc model.SendCoinRequest

	if err := jsoniter.NewDecoder(r.Body).Decode(&sc); err != nil {
		h.Logger.Sugar.Infow("error in request handler", "error: ", err)
		SendErrors(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(model.ValidLogin{}).(*model.Claims).UserdID

	err := h.Stor.SendCoin(r.Context(), userID, sc)
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
		// switch err {
		// case apperr.ErrInsufficientFunds:
		// 	SendErrors(w, err.Error(), http.StatusBadRequest)
		// case apperr.ErrSenderNotFound:
		// 	SendErrors(w, err.Error(), http.StatusBadRequest)
		// case apperr.ErrRecipientNotFound:
		// 	SendErrors(w, err.Error(), http.StatusBadRequest)
		// default:
		// 	SendErrors(w, err.Error(), http.StatusInternalServerError)
		// }

		// h.Logger.Sugar.Infow("send coin error", "error: ", err)
		// return
	}
	w.WriteHeader(http.StatusOK)
}
