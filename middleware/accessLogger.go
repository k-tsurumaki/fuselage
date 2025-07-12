package middleware

import (
	"log"
	"time"

	"github.com/k-tsurumaki/fuselage"
)

type LoggerConfig struct {
	// Skipper defines a function to skip middleware
	Skipper func(*fuselage.Context) bool
}

var DefaultLoggerConfig = LoggerConfig{
	Skipper: func(c *fuselage.Context) bool {
		return false
	},
}

func Logger() fuselage.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

func LoggerWithConfig(config LoggerConfig) fuselage.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultLoggerConfig.Skipper
	}

	return func(next fuselage.HandlerFunc) fuselage.HandlerFunc {
		return func(c *fuselage.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			start := time.Now()
			requestID := GetRequestID(c)

			err := next(c)

			duration := time.Since(start)
			status := c.Status()
			if !c.IsWritten() && status == 0 {
				status = 200
			}

			log.Printf("[%s] %s %s %d %v", requestID, c.Request.Method, c.Request.URL.Path, status, duration)
			return err
		}
	}
}
