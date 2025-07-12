package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/k-tsurumaki/fuselage"
)

type RequestIDConfig struct {
	// Generator defines a function to generate an ID.
	Generator func() string

	// RequestIDHandler defines a function which is executed for a request id.
	RequestIDHandler func(fuselage.Context, string)

	// TargetHeader defines what header to look for to populate the id
	TargetHeader string

	// Skipper defines a function to skip middleware
	Skipper func(*fuselage.Context) bool
}

var DefaultRequestIDConfig = RequestIDConfig{
	Generator:    generateRequestID,
	TargetHeader: fuselage.HeaderXRequestID,
	Skipper: func(c *fuselage.Context) bool {
		return false
	},
}

func RequestID() fuselage.MiddlewareFunc {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestID middleware adds unique request ID
func RequestIDWithConfig(config RequestIDConfig) fuselage.MiddlewareFunc {
	if config.TargetHeader == "" {
		config.TargetHeader = fuselage.HeaderXRequestID
	}
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIDConfig.Skipper
	}

	return func(next fuselage.HandlerFunc) fuselage.HandlerFunc {
		return func(c *fuselage.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			rid := c.Header(config.TargetHeader)

			// If request ID is not provided, generate a new one
			if rid == "" {
				rid = config.Generator()
			}

			// Set request ID in context and response header
			ctx := context.WithValue(c.Request.Context(), fuselage.RequestIDKey, rid)
			c.Request = c.Request.WithContext(ctx)
			c.SetHeader(config.TargetHeader, rid)

			if config.RequestIDHandler != nil {
				config.RequestIDHandler(*c, rid)
			}

			return next(c)
		}
	}
}

// GetRequestID extracts request ID from context
func GetRequestID(c *fuselage.Context) string {
	if id := c.Request.Context().Value(fuselage.RequestIDKey); id != nil {
		return id.(string)
	}
	return "unknown"
}

func generateRequestID() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
