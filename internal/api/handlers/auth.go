package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgconn"
	jsoniter "github.com/json-iterator/go"
	"github.com/plasmatrip/avito_merch/internal/apperr"
	"github.com/plasmatrip/avito_merch/internal/model"
	"github.com/rgurov/pgerrors"
)

func (h *Handlers) Auth(w http.ResponseWriter, r *http.Request) {
	var req model.AuthRequest

	if err := jsoniter.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Sugar.Infow("error in request handler", "error: ", err)
		SendErrors(w, apperr.ErrBadJSON)
		return
	}

	if len(req.UserName) == 0 || len(req.Password) == 0 {
		h.Logger.Sugar.Infow("error in authentication data", "error: ", errors.New("empty login or password"))
		SendErrors(w, apperr.ErrEmptyLoginOrPassword)
		return
	}

	id, err := h.Stor.UserAuth(r.Context(), req)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrors.UniqueViolation {
				h.Logger.Sugar.Infow("authentication error", "error: ", err)
				SendErrors(w, apperr.ErrAuthenticationError)
				return
			}
		}

		h.Logger.Sugar.Infow("internal error", "error: ", err)
		SendErrors(w, err)
		return
	}

	token, err := h.Token(id, req)
	if err != nil {
		h.Logger.Sugar.Infow("error generating JWT", "error: ", err)
		SendErrors(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(model.AuthResponse{Token: token}); err != nil {
		h.Logger.Sugar.Infow("error encoding response", "error: ", err)
		SendErrors(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) Token(id uuid.UUID, lr model.AuthRequest) (string, error) {
	claims := model.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Subject:   lr.UserName,
		},
		UserdID: id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(h.Config.TokenSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
