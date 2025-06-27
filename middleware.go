package fuselage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"
)

// Logger middleware logs HTTP requests with RequestID
func Logger(next HandlerFunc) HandlerFunc {
	return func(c *Context) error {
		start := time.Now()
		requestID := GetRequestID(c)

		err := next(c)

		duration := time.Since(start)
		status := c.status
		if status == 0 {
			status = 200
		}

		log.Printf("[%s] %s %s %d %v", requestID, c.Request.Method, c.Request.URL.Path, status, duration)
		return err
	}
}

// Recover middleware recovers from panics
func Recover(next HandlerFunc) HandlerFunc {
	return func(c *Context) error {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c)
				log.Printf("[%s] panic recovered: %v", requestID, err)
				_ = c.String(500, "Internal Server Error")
			}
		}()
		return next(c)
	}
}

// Timeout middleware adds request timeout
func Timeout(duration time.Duration) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
			defer cancel()

			c.Request = c.Request.WithContext(ctx)
			return next(c)
		}
	}
}

// RequestID middleware adds unique request ID
func RequestID(next HandlerFunc) HandlerFunc {
	return func(c *Context) error {
		requestID := generateRequestID()
		ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Request-ID", requestID)
		return next(c)
	}
}

// GetRequestID extracts request ID from context
func GetRequestID(c *Context) string {
	if id := c.Request.Context().Value(RequestIDKey); id != nil {
		return id.(string)
	}
	return "unknown"
}

func generateRequestID() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
