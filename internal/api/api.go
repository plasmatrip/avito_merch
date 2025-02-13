package api

import "net/http"

type API interface {
	Info(w http.ResponseWriter, r *http.Request)
	Buy(w http.ResponseWriter, r *http.Request)
	SendCoin(w http.ResponseWriter, r *http.Request)
	Auth(w http.ResponseWriter, r *http.Request)
}
