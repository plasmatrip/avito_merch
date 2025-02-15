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
		SendErrors(w, err)
		return
	}

	if len(sc.ToUser) == 0 && sc.Amount == 0 {
		h.Logger.Sugar.Infow("error in request handler", "error: ", apperr.ErrBadJSON)
		SendErrors(w, apperr.ErrBadJSON)
		return
	}

	if sc.Amount <= 0 {
		h.Logger.Sugar.Infow("error in request handler", "error: ", apperr.ErrAmonutIsLessThanOrEqualToZero)
		SendErrors(w, apperr.ErrAmonutIsLessThanOrEqualToZero)
		return
	}

	userID := r.Context().Value(model.ValidLogin{}).(*model.Claims).UserdID

	err := h.Stor.SendCoin(r.Context(), userID, sc)
	if err != nil {
		SendErrors(w, err)
		h.Logger.Sugar.Infow("internal error", "error: ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
