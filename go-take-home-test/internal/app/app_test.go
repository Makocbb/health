package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-take-home-test/internal/app"
)

func TestPostIngest(t *testing.T) {
	e := app.New()

	req := httptest.NewRequest(http.MethodPost, "/ingest", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
