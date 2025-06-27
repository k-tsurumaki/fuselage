package fuselage

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRouter_GET(t *testing.T) {
	router := New()
	router.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET test"))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "GET test" {
		t.Errorf("Expected body 'GET test', got '%s'", w.Body.String())
	}
}

func TestRouter_POST(t *testing.T) {
	router := New()
	router.POST("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("POST test"))
	})

	req := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestRouter_ParameterExtraction(t *testing.T) {
	router := New()
	router.GET("/users/:id", func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r, "id")
		w.Write([]byte("User ID: " + id))
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	expected := "User ID: 123"
	if w.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, w.Body.String())
	}
}

func TestRouter_NotFound(t *testing.T) {
	router := New()

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestMiddleware_Logger(t *testing.T) {
	router := New()
	router.Use(Logger)
	router.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestMiddleware_Recover(t *testing.T) {
	router := New()
	router.Use(Recover)
	router.GET("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestMiddleware_Timeout(t *testing.T) {
	router := New()
	router.Use(Timeout(50 * time.Millisecond))
	router.GET("/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Err() != nil {
			t.Log("Context cancelled as expected")
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Test passes if no panic occurs
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}