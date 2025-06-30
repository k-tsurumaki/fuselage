package middleware

import (
	"context"
	"time"

	"github.com/k-tsurumaki/fuselage"
)

type TimeoutConfig struct {
	// Duration defines the timeout duration for requests.
	Timeout time.Duration
}

func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		Timeout: 30 * time.Second, // Default timeout duration
	}
}

func Timeout() fuselage.MiddlewareFunc {
	return TimeoutWithConfig(DefaultTimeoutConfig())
}

// Timeout middleware adds request timeout
func TimeoutWithConfig(config TimeoutConfig) fuselage.MiddlewareFunc {
	return func(next fuselage.HandlerFunc) fuselage.HandlerFunc {
		return func(c *fuselage.Context) error {
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
