package testserver

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers_RoundTrip(t *testing.T) {
	handlers := NewHandlers(
		Handler{Status: http.StatusCreated, Response: []byte("success")},
		Handler{Response: []byte("OK")},
		Handler{Status: http.StatusConflict},
		Handler{Path: "/rest", Status: http.StatusCreated, Response: []byte("rested"), Headers: map[string]string{"Content-Type": "application/json"}},
		Handler{Path: "rest", Status: http.StatusTeapot, Response: []byte("teapot")},
	)

	tests := []struct {
		status int
		body   string
		path   string
	}{
		{status: 201, body: "success", path: "/"},
		{status: 201, body: "rested", path: "/rest"},
		{status: 200, body: "OK", path: "/"},
		{status: 409, body: "", path: "/"},
		{status: 418, body: "teapot", path: "/rest"},
		{status: 201, body: "success", path: "/"},
		{status: 201, body: "rested", path: "/rest"},
	}
	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		resp, _ := handlers.RoundTrip(req)
		if resp.StatusCode != tt.status {
			t.Errorf("Handlers.RoundTrip().Status = %v, want %v", resp.StatusCode, tt.status)
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Unexpected response body failure = %v", err)
		}
		if string(b) != tt.body {
			t.Errorf("Handlers.RoundTrip().Body = %v, want %v", string(b), tt.body)
		}
		resp.Body.Close()
	}
}

func TestHandlers_Request_empty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := NewHandlers().Request(req)
	if err == nil {
		t.Errorf("expected server response failure = %v", err)
		return
	}
	if err.Error() != "no path available" {
		t.Errorf(" NewHandlers().Request().Error = %v, want %v", err.Error(), "no path available")
	}
}
