package middleware

import (
	"net/http"

	"golang.org/x/time/rate"

	"github.com/plasmatrip/avito_merch/internal/logger"
)

func WithLimitter(log logger.Logger) func(next http.Handler) http.Handler {
	var limiter = rate.NewLimiter(rate.Limit(1100), 200)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				log.Sugar.Infow("too many requests", "request", r.RequestURI)
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
