package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChiServer(t *testing.T) {
	ctx := context.Background()
	chiServer := NewChiServer(ctx)
	chiServer.MountRoutes()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	chiServer.Handler().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"ok"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
