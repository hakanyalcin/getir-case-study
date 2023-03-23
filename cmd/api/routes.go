package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Post("/records", app.getRecords)
	mux.Get("/in-memory", app.getEntry)
	mux.Post("/in-memory", app.setEntry)

	return mux
}
