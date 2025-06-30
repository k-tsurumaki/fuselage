package middleware

import (
	"log"

	"github.com/k-tsurumaki/fuselage"
)

type RecoverConfig struct {}

var DefaultRecoverConfig = RecoverConfig{}

func Recover() fuselage.MiddlewareFunc {
	return RecoverWithConfig(DefaultRecoverConfig)
}

// Recover middleware recovers from panics
func RecoverWithConfig(config RecoverConfig) fuselage.MiddlewareFunc {
	return func(next fuselage.HandlerFunc) fuselage.HandlerFunc {
		return func(c *fuselage.Context) error {
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
}
