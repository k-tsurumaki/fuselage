package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k-tsurumaki/fuselage"
)

func TestRequestID(t *testing.T) {
	router := fuselage.New()
	router.Use(RequestID())
	
	router.GET("/test", func(c *fuselage.Context) error {
		id := GetRequestID(c)
		if id == "" || id == "unknown" {
			t.Errorf("Expected request ID to be set")
		}
		return c.String(http.StatusOK, id)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	
	if rec.Header().Get("X-Request-ID") == "" {
		t.Errorf("Expected X-Request-ID header to be set")
	}
}

func TestRequestIDWithCustomGenerator(t *testing.T) {
	router := fuselage.New()
	router.Use(RequestIDWithConfig(RequestIDConfig{
		Generator: func() string {
			return "custom-id"
		},
	}))
	
	router.GET("/test", func(c *fuselage.Context) error {
		id := GetRequestID(c)
		if id != "custom-id" {
			t.Errorf("Expected custom-id, got %s", id)
		}
		return c.String(http.StatusOK, id)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	
	if rec.Body.String() != "custom-id" {
		t.Errorf("Expected custom-id in response")
	}
}

func TestRequestIDWithSkipper(t *testing.T) {
	router := fuselage.New()
	router.Use(RequestIDWithConfig(RequestIDConfig{
		Skipper: func(c *fuselage.Context) bool {
			return c.Request.URL.Path == "/skip"
		},
	}))
	
	router.GET("/skip", func(c *fuselage.Context) error {
		return c.String(http.StatusOK, "Skipped")
	})

	req := httptest.NewRequest("GET", "/skip", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}