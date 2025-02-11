package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/plasmatrip/avito_merch/internal/api/handlers"
	"github.com/plasmatrip/avito_merch/internal/logger"
	"github.com/plasmatrip/avito_merch/internal/model"
)

func WithAuthentication(log logger.Logger, tokenSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Sugar.Info("missing authorization header")
				handlers.SendErrors(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Sugar.Infow("invalid authorization header format", "parts", parts)
				handlers.SendErrors(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			claims := &model.Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims,
				func(t *jwt.Token) (interface{}, error) {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
					}
					return []byte(tokenSecret), nil
				})

			if err != nil {
				log.Sugar.Infow("JWT token error", "error", err)
				handlers.SendErrors(w, "JWT token error", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				log.Sugar.Infow("invalid token", "token", token)
				handlers.SendErrors(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), model.ValidLogin{}, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
