package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestCreateGameHandler(t *testing.T) {
	t.Run("should reject GET requests", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/create", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		initInstanceManager()
		handler := http.HandlerFunc(createGameHandler)

		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %v", status)
		}

		// Check the response body is what we expect.
		expected := "method not permitted\n"
		if rr.Body.String() != expected {
			t.Errorf("expected %v from handler, got %v", expected, rr.Body.String())
		}
	})

	t.Run("should create games and return their ID", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/create", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		initInstanceManager()
		handler := http.HandlerFunc(createGameHandler)

		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("expected status 200, got %v", status)
		}

		// Check the response body is what we expect.
		expected := `{"id":".*"}`
		match, _ := regexp.MatchString(expected, rr.Body.String())
		if !match {
			t.Errorf("expected body to match %v, got %v", expected, rr.Body.String())
		}

		if len(im.Instances) != 1 {
			t.Errorf("expected an instance to have been created")
		}
	})
}
