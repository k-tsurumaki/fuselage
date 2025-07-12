package middleware

import (
	"log"

	"github.com/k-tsurumaki/fuselage"
)

type RecoverConfig struct {
	// Skipper defines a function to skip middleware
	Skipper func(*fuselage.Context) bool
	// ErrorHandler defines a function to handle panic errors
	ErrorHandler func(*fuselage.Context, interface{}) error
}

var DefaultRecoverConfig = RecoverConfig{
	Skipper: func(c *fuselage.Context) bool {
		return false
	},
	ErrorHandler: func(c *fuselage.Context, err interface{}) error {
		requestID := GetRequestID(c)
		log.Printf("[%s] panic recovered: %v", requestID, err)
		return c.String(500, "Internal Server Error")
	},
}

func Recover() fuselage.MiddlewareFunc {
	return RecoverWithConfig(DefaultRecoverConfig)
}

// Recover middleware recovers from panics
func RecoverWithConfig(config RecoverConfig) fuselage.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultRecoverConfig.Skipper
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = DefaultRecoverConfig.ErrorHandler
	}

	return func(next fuselage.HandlerFunc) fuselage.HandlerFunc {
		return func(c *fuselage.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if err := recover(); err != nil {
					_ = config.ErrorHandler(c, err)
				}
			}()
			return next(c)
		}
	}
}
