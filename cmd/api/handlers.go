package main

import (
	"encoding/json"
	"getir-case-study/internal/models"

	"net/http"
)

func (app *application) getRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.unsupportedMethod(w, r)
		return

	}
	var payload models.RecordPayload
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	record := models.RecordPayload{
		StartDate: payload.StartDate,
		EndDate:   payload.EndDate,
		MinCount:  payload.MinCount,
		MaxCount:  payload.MaxCount,
	}

	res, err := app.DB.GetRecords(record)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, res)
}

func (app *application) setEntry(w http.ResponseWriter, r *http.Request) {
	var payload models.CachePayload
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	entry := models.CachePayload{
		Key:   payload.Key,
		Value: payload.Value,
	}
	res, err := app.cache.SetEntry(entry)
	if err != nil {
		return
	}
	out, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) getEntry(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("key")

	res, err := app.cache.GetEntry(id)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	out, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)

}
