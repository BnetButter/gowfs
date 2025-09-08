package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

)

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil);
	rec := httptest.NewRecorder()
	healthCheckHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200; got %d", res.StatusCode)
	}

	body := rec.Body.String()
	if body != "OK" {
		t.Errorf("Unexpected body: %s", body);
	}
}

