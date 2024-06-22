package controllers

import (
	"net/http"
)

func (app *App) Healthz(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}