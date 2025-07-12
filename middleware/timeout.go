package middleware

import (
	"context"
	"time"

	"github.com/k-tsurumaki/fuselage"
)

type TimeoutConfig struct {
	// Duration defines the timeout duration for requests.
	Timeout time.Duration
	// Skipper defines a function to skip middleware
	Skipper func(*fuselage.Context) bool
	// ErrorHandler defines a function to handle timeout errors
	ErrorHandler func(*fuselage.Context) error
}

func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		Timeout: 30 * time.Second, // Default timeout duration
		Skipper: func(c *fuselage.Context) bool {
			return false
		},
		ErrorHandler: func(c *fuselage.Context) error {
			return c.String(408, "Request Timeout")
		},
	}
}

func Timeout() fuselage.MiddlewareFunc {
	return TimeoutWithConfig(DefaultTimeoutConfig())
}

// Timeout middleware adds request timeout
func TimeoutWithConfig(config TimeoutConfig) fuselage.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultTimeoutConfig().Skipper
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = DefaultTimeoutConfig().ErrorHandler
	}

	return func(next fuselage.HandlerFunc) fuselage.HandlerFunc {
		return func(c *fuselage.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			if config.Timeout <= 0 {
				return next(c) // No timeout, just proceed
			}
			
			ctx, cancel := context.WithTimeout(c.Request.Context(), config.Timeout)
			defer cancel()

			c.Request = c.Request.WithContext(ctx)
			return next(c)
		}
	}
}
