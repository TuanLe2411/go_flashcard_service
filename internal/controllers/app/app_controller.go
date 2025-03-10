package app

import "net/http"

type AppController struct{}

func (a *AppController) HeathCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
