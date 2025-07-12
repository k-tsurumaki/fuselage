package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k-tsurumaki/fuselage"
)

func TestLogger(t *testing.T) {
	router := fuselage.New()
	router.Use(Logger())
	
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

func TestLoggerWithSkipper(t *testing.T) {
	router := fuselage.New()
	router.Use(LoggerWithConfig(LoggerConfig{
		Skipper: func(c *fuselage.Context) bool {
			return c.Request.URL.Path == "/skip"
		},
	}))
	
	router.GET("/skip", func(c *fuselage.Context) error {
		return c.String(http.StatusOK, "Skipped")
	})
	router.GET("/log", func(c *fuselage.Context) error {
		return c.String(http.StatusOK, "Logged")
	})

	// Test skipped path
	req1 := httptest.NewRequest("GET", "/skip", nil)
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	
	if rec1.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec1.Code)
	}

	// Test logged path
	req2 := httptest.NewRequest("GET", "/log", nil)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	
	if rec2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec2.Code)
	}
}