package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k-tsurumaki/fuselage"
)

func TestCORS(t *testing.T) {
	router := fuselage.New()
	router.Use(CORS())
	
	router.GET("/test", func(c *fuselage.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	
	if rec.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("Expected CORS origin header")
	}
}

func TestCORSPreflight(t *testing.T) {
	router := fuselage.New()
	router.Use(CORS())
	
	router.GET("/test", func(c *fuselage.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", rec.Code)
	}
}

func TestCORSWithSkipper(t *testing.T) {
	router := fuselage.New()
	router.Use(CORSWithConfig(&CORSConfig{
		AllowOrigins: []string{"*"},
		Skipper: func(c *fuselage.Context) bool {
			return c.Request.URL.Path == "/skip"
		},
	}))
	
	router.GET("/skip", func(c *fuselage.Context) error {
		return c.String(http.StatusOK, "Skipped")
	})

	req := httptest.NewRequest("GET", "/skip", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}