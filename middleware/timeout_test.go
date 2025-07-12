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
		Timeout: 10 * time.Millisecond,
	}))
	
	router.GET("/fast", func(c *fuselage.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	
	router.GET("/slow", func(c *fuselage.Context) error {
		time.Sleep(20 * time.Millisecond)
		return c.String(http.StatusOK, "OK")
	})

	// Test fast request
	req1 := httptest.NewRequest("GET", "/fast", nil)
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	
	if rec1.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec1.Code)
	}

	// Test slow request (should timeout)
	req2 := httptest.NewRequest("GET", "/slow", nil)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	
	if rec2.Code != http.StatusRequestTimeout {
		t.Errorf("Expected status 408, got %d", rec2.Code)
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
		time.Sleep(20 * time.Millisecond)
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
		Timeout: 10 * time.Millisecond,
		ErrorHandler: func(c *fuselage.Context) error {
			return c.String(http.StatusServiceUnavailable, "Custom timeout")
		},
	}))
	
	router.GET("/slow", func(c *fuselage.Context) error {
		time.Sleep(20 * time.Millisecond)
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/slow", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", rec.Code)
	}
	
	if rec.Body.String() != "Custom timeout" {
		t.Errorf("Expected custom timeout message")
	}
}