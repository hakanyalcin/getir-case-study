package main

import (
	"encoding/json"
	"errors"
	"getir-case-study/internal/models"
	"io"
	"net/http"
)

// create custom success response
func (app *application) successResponse(w http.ResponseWriter, status int, data []models.Record) error {
	response := models.RecordResponse{
		Code:    0,
		Msg:     "success",
		Records: data,
	}
	out, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)

	return nil
}

// readBody reads json from request body into data. We only accept a single json value in the body
func (app *application) readBody(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // max one megabyte in request body
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	// we only allow one entry in the json file
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("invalid body: Body must only have a single JSON value")
	}

	return nil
}

// describing the error message in the response
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) error {
	var payload struct {
		Code    int    `json:"error"`
		Message string `json:"message"`
	}

	payload.Code = 1
	payload.Message = err.Error()

	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(out)
	return nil
}

func (app *application) unsupportedMethodResponse(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Code    int    `json:"error"`
		Message string `json:"message"`
	}

	payload.Code = 2
	payload.Message = "MethodNotAllowed: Unsupported Method"

	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(out)
}

// if key not found in the cache, then return custom message
func (app *application) cacheMissingResponse(w http.ResponseWriter, r *http.Request, err error) error {
	var payload struct {
		Message string `json:"message"`
	}

	payload.Message = err.Error()

	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(out)
	return nil
}

// error occurs when fetching from database or localcache
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) error {
	var payload struct {
		Code    int    `json:"error"`
		Message string `json:"message"`
	}

	payload.Code = 3
	payload.Message = err.Error()

	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(out)
	return nil
}
