package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k-tsurumaki/fuselage"
)

func TestRecover(t *testing.T) {
	router := fuselage.New()
	router.Use(Recover())
	
	router.GET("/panic", func(c *fuselage.Context) error {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rec.Code)
	}
}

func TestRecoverWithCustomErrorHandler(t *testing.T) {
	router := fuselage.New()
	router.Use(RecoverWithConfig(RecoverConfig{
		ErrorHandler: func(c *fuselage.Context, err interface{}) error {
			return c.String(http.StatusBadRequest, "Custom error")
		},
	}))
	
	router.GET("/panic", func(c *fuselage.Context) error {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
	
	if rec.Body.String() != "Custom error" {
		t.Errorf("Expected custom error message")
	}
}

func TestRecoverWithSkipper(t *testing.T) {
	router := fuselage.New()
	router.Use(RecoverWithConfig(RecoverConfig{
		Skipper: func(c *fuselage.Context) bool {
			return c.Request.URL.Path == "/skip"
		},
	}))
	
	router.GET("/skip", func(c *fuselage.Context) error {
		panic("should not be recovered")
	})

	req := httptest.NewRequest("GET", "/skip", nil)
	rec := httptest.NewRecorder()
	
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic to not be recovered")
		}
	}()
	
	router.ServeHTTP(rec, req)
}