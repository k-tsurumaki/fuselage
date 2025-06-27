package fuselage

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter_GET(t *testing.T) {
	router := New()
	_ = router.GET("/test", func(c *Context) error {
		return c.String(http.StatusOK, "GET test")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "GET test" {
		t.Errorf("Expected body 'GET test', got '%s'", w.Body.String())
	}
}

func TestRouter_DuplicateRoute(t *testing.T) {
	router := New()
	err1 := router.GET("/test", func(c *Context) error { return nil })
	err2 := router.GET("/test", func(c *Context) error { return nil })

	if err1 != nil {
		t.Errorf("First route registration should succeed")
	}
	if err2 == nil {
		t.Errorf("Duplicate route registration should fail")
	}
}

func TestRouter_ParameterExtraction(t *testing.T) {
	router := New()
	_ = router.GET("/users/:id", func(c *Context) error {
		id := c.Param("id")
		return c.String(http.StatusOK, "User ID: "+id)
	})

	req := httptest.NewRequest("GET", "/users/123", http.NoBody)
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

	req := httptest.NewRequest("GET", "/nonexistent", http.NoBody)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestRouter_MethodNotAllowed(t *testing.T) {
	router := New()
	_ = router.GET("/test", func(c *Context) error { return nil })

	req := httptest.NewRequest("POST", "/test", http.NoBody)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestContext_JSON(t *testing.T) {
	router := New()
	_ = router.GET("/json", func(c *Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "hello"})
	})

	req := httptest.NewRequest("GET", "/json", http.NoBody)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type should be application/json")
	}

	if !strings.Contains(w.Body.String(), "hello") {
		t.Error("Response should contain 'hello'")
	}
}

func TestContext_ParamInt(t *testing.T) {
	router := New()
	_ = router.GET("/users/:id", func(c *Context) error {
		id, err := c.ParamInt("id")
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID")
		}
		return c.JSON(http.StatusOK, map[string]int{"id": id})
	})

	req := httptest.NewRequest("GET", "/users/123", http.NoBody)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Body.String(), "123") {
		t.Error("Response should contain '123'")
	}
}

func TestMiddleware_RequestID(t *testing.T) {
	router := New()
	router.Use(RequestID)
	_ = router.GET("/test", func(c *Context) error {
		requestID := GetRequestID(c)
		if requestID == "unknown" {
			t.Error("RequestID should be set")
		}
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") == "" {
		t.Error("X-Request-ID header should be set")
	}
}

func TestValidateStruct(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name" validate:"required,min=2"`
	}

	// Valid struct
	valid := TestStruct{Name: "John"}
	errors := ValidateStruct(valid)
	if len(errors) > 0 {
		t.Errorf("Valid struct should not have errors: %v", errors)
	}

	// Invalid struct - missing required field
	invalid := TestStruct{}
	errors = ValidateStruct(invalid)
	if len(errors) == 0 {
		t.Error("Invalid struct should have errors")
	}
}