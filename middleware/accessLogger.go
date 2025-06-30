package middleware

import (
	"log"
	"time"

	"github.com/k-tsurumaki/fuselage"
)

type LoggerConfig struct {}

var DefaultLoggerVonfig = LoggerConfig{}

func Logger() fuselage.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerVonfig)
}

func LoggerWithConfig(config LoggerConfig) fuselage.MiddlewareFunc {
	return func(next fuselage.HandlerFunc) fuselage.HandlerFunc {
		return func(c *fuselage.Context) error {
			start := time.Now()
			requestID := GetRequestID(c)

			err := next(c)

			duration := time.Since(start)
			status := c.Status()
			if status == 0 {
				status = 200
			}

			log.Printf("[%s] %s %s %d %v", requestID, c.Request.Method, c.Request.URL.Path, status, duration)
			return err
		}
	}
}
