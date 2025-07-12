package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/k-tsurumaki/fuselage"
)

func TestRateLimit(t *testing.T) {
	router := fuselage.New()
	
	// Configure rate limit: 2 requests per second
	router.Use(RateLimitWithConfig(RateLimitConfig{
		Limit:  2,
		Window: time.Second,
	}))
	
	router.GET("/test", func(c *fuselage.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// First request should pass
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "127.0.0.1:8080"
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	
	if rec1.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec1.Code)
	}

	// Second request should pass
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "127.0.0.1:8080"
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	
	if rec2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec2.Code)
	}

	// Third request should be rate limited
	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "127.0.0.1:8080"
	rec3 := httptest.NewRecorder()
	router.ServeHTTP(rec3, req3)
	
	if rec3.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", rec3.Code)
	}
}

func TestRateLimitDifferentIPs(t *testing.T) {
	router := fuselage.New()
	
	router.Use(RateLimitWithConfig(RateLimitConfig{
		Limit:  1,
		Window: time.Second,
	}))
	
	router.GET("/test", func(c *fuselage.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Request from IP 1
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "127.0.0.1:8080"
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	
	if rec1.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec1.Code)
	}

	// Request from IP 2 should still pass
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.1:8080"
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	
	if rec2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec2.Code)
	}
}