package main

import (
	"bytes"
	"context"
	"encoding/json"
	"getir-case-study/internal/cache"
	"getir-case-study/internal/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestApiForSetEntry(t *testing.T) {
	app := &application{
		cache: cache.NewLocalCache(),
	}
	// Create a new test server
	ts := httptest.NewServer(http.HandlerFunc(app.setEntry))
	defer ts.Close()

	// Create a new entry
	newEntry := models.CachePayload{
		Key:   "test_key",
		Value: "test_value",
	}

	// Convert entry to JSON
	jsonEntry, err := json.Marshal(newEntry)
	if err != nil {
		t.Fatal(err)
	}

	// Send POST request to set entry
	resp, err := http.Post(ts.URL+"/in-memory", "application/json", bytes.NewBuffer(jsonEntry))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200; got %d", resp.StatusCode)
	}

	var res cache.Entry
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		t.Fatal(err)
	}

	// Check response values
	if res.Key != newEntry.Key {
		t.Errorf("Expected key %s; got %s", newEntry.Key, res.Key)
	}
	if res.Value != newEntry.Value {
		t.Errorf("Expected value %s; got %s", newEntry.Value, res.Value)
	}
}

func TestApiForGetEntryCacheMissing(t *testing.T) {
	app := &application{
		cache: cache.NewLocalCache(),
	}
	req, err := http.NewRequest("GET", "/in-memory?key=test_key", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.getEntry)
	handler.ServeHTTP(rr, req)
	cacheMissingMessage := "{\n\t\"message\": \"cache missing: The key isn't in cache\"\n}"
	if rr.Body.String() != cacheMissingMessage {
		t.Errorf(" returned wrong message: got %v want %v",
			rr.Body.String(), cacheMissingMessage)
	}
}

func TestApiForGetEntry(t *testing.T) {
	app := &application{
		cache: cache.NewLocalCache(),
	}
	ts := httptest.NewServer(http.HandlerFunc(app.setEntry))
	defer ts.Close()

	// Create a new entry
	newEntry := models.CachePayload{
		Key:   "test_key",
		Value: "test_value",
	}

	// Convert entry to JSON
	jsonEntry, err := json.Marshal(newEntry)
	if err != nil {
		t.Fatal(err)
	}

	// Send POST request to set entry
	setResp, err := http.Post(ts.URL+"/in-memory", "application/json", bytes.NewBuffer(jsonEntry))
	if err != nil {
		t.Fatal(err)
	}
	defer setResp.Body.Close()

	req, err := http.NewRequest("GET", "/in-memory?key=test_key", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.getEntry)
	handler.ServeHTTP(rr, req)
	expectedBody := "{\n\t\"key\": \"test_key\",\n\t\"value\": \"test_value\"\n}"
	if rr.Body.String() != expectedBody {
		t.Errorf(" returned wrong message: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}
func TestGetRecords(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getircase-study?retryWrites=true"))
	require.NoError(t, err)

	app := &application{
		DB: models.DBModel{DB: client},
	}

	req := httptest.NewRequest(http.MethodPost, "/records", nil)

	// Create a mocked http.ResponseWriter and http.Request with the necessary request headers and payload
	rr := httptest.NewRecorder()
	payload := `{"startDate":"2016-01-26","endDate":"2018-02-02","minCount":2700,"maxCount":3000}`
	req.Body = ioutil.NopCloser(strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	// Call the getRecords function with the mocked http.ResponseWriter and http.Request
	app.getRecords(rr, req)

	// Verify that the correct response code and payload are returned
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"key":"TAKwGc6Jr4i8Z487","createdAt":"2017-01-28T01:22:14.398Z","totalCount":2800},
	{"key":"NAeQ8eX7e5TEg7oH","createdAt":"2017-01-27T08:19:14.135Z"","totalCount":2900}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
