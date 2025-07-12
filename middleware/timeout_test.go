package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/k-tsurumaki/fuselage"
)

func TestTimeout(t *testing.T) {
	router := fuselage.New()
	router.Use(TimeoutWithConfig(TimeoutConfig{
		Timeout: 100 * time.Millisecond,
	}))
	
	router.GET("/test", func(c *fuselage.Context) error {
		// Check if context has timeout
		if c.Request.Context().Err() != nil {
			return c.String(http.StatusRequestTimeout, "Timeout")
		}
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestTimeoutWithSkipper(t *testing.T) {
	router := fuselage.New()
	router.Use(TimeoutWithConfig(TimeoutConfig{
		Timeout: 10 * time.Millisecond,
		Skipper: func(c *fuselage.Context) bool {
			return c.Request.URL.Path == "/skip"
		},
	}))
	
	router.GET("/skip", func(c *fuselage.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/skip", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestTimeoutWithCustomErrorHandler(t *testing.T) {
	router := fuselage.New()
	router.Use(TimeoutWithConfig(TimeoutConfig{
		Timeout: 100 * time.Millisecond,
		ErrorHandler: func(c *fuselage.Context) error {
			return c.String(http.StatusServiceUnavailable, "Custom timeout")
		},
	}))
	
	router.GET("/test", func(c *fuselage.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}