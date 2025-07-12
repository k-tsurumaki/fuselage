package fuselage

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContext_IsWritten(t *testing.T) {
	// Test initial state
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	c := &Context{
		Request:  req,
		Response: w,
		written:  false,
	}

	if c.IsWritten() {
		t.Error("Expected IsWritten() to be false initially")
	}

	// Test after JSON response
	c.JSON(http.StatusOK, map[string]string{"test": "data"})
	if !c.IsWritten() {
		t.Error("Expected IsWritten() to be true after JSON response")
	}
}

func TestContext_IsWritten_String(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	c := &Context{
		Request:  req,
		Response: w,
		written:  false,
	}

	// Test after String response
	c.String(http.StatusOK, "test")
	if !c.IsWritten() {
		t.Error("Expected IsWritten() to be true after String response")
	}
}

func TestContext_IsWritten_SetStatus(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	c := &Context{
		Request:  req,
		Response: w,
		written:  false,
	}

	// Test after SetStatus
	c.SetStatus(http.StatusCreated)
	if !c.IsWritten() {
		t.Error("Expected IsWritten() to be true after SetStatus")
	}
}