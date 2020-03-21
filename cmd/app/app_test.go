package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrzacarias/private-kit-server/internal/mock"
)

func init() {
	InfectedClient = &mock.InfectedClient{}
}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Using ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestInfectedHandler(test *testing.T) {
	test.Run("Infected request works", func(t *testing.T) {
		var jsonStr = []byte(`{"uuids":["ef872b78-3df5-483c-a764-6f33e1a23898"],"since_ts":"2020-03-21T16:36:59+00:00"}`)
		req, err := http.NewRequest("POST", "/infected", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}

		// Using ResponseRecorder to record the response
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(InfectedHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	test.Run("Infected request do not work", func(t *testing.T) {
		var jsonStr = []byte(`{"uuids":["ef872b78-3df5-483c-a764-6f33e1a23898"]`)
		req, err := http.NewRequest("POST", "/infected", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}

		// Using ResponseRecorder to record the response
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(InfectedHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}
