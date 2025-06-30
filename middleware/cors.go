package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/k-tsurumaki/fuselage"
)

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

var DefaultCORSConfig = CORSConfig{
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{fuselage.GET, fuselage.HEAD, fuselage.PUT, fuselage.PATCH, fuselage.POST, fuselage.DELETE},
}

func CORS() fuselage.MiddlewareFunc {
	return CORSWithConfig(&DefaultCORSConfig)
}

func CORSWithConfig(config *CORSConfig) fuselage.MiddlewareFunc {
	if len(config.AllowedOrigins) == 0 {
		config.AllowedOrigins = DefaultCORSConfig.AllowedOrigins
	}
	if len(config.AllowedMethods) == 0 {
		config.AllowedMethods = DefaultCORSConfig.AllowedMethods
	}
	if len(config.AllowedHeaders) == 0 {
		config.AllowedHeaders = []string{fuselage.HeaderContentType, fuselage.HeaderAuthorization}
	}

	allowedOriginPatterns := make([]*regexp.Regexp, 0, len(config.AllowedOrigins))
	for _, origin := range config.AllowedOrigins {
		if origin == "*" {
			continue
		}
		pattern := regexp.QuoteMeta(origin)
		pattern = strings.ReplaceAll(pattern, "\\*", ".*")
		pattern = strings.ReplaceAll(pattern, "\\?", ".")
		pattern = "^" + pattern + "$"
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		allowedOriginPatterns = append(allowedOriginPatterns, re)
	}

	allowedMethods := strings.Join(config.AllowedMethods, ", ")
	allowedHeaders := strings.Join(config.AllowedHeaders, ", ")

	return func(next fuselage.HandlerFunc) fuselage.HandlerFunc {
		return func(c *fuselage.Context) error {
			origin := c.Request.Header.Get(fuselage.HeaderOrigin)
			allowed := false
			if origin != "" {
				for _, o := range config.AllowedOrigins {
					if o == "*" || o == origin {
						allowed = true
						break
					}
				}
				if !allowed {
					for _, re := range allowedOriginPatterns {
						if re.MatchString(origin) {
							allowed = true
							break
						}
					}
				}
				if allowed {
					c.SetHeader(fuselage.HeaderAccessControlAllowOrigin, origin)
				}
			}

			c.SetHeader(fuselage.HeaderAccessControlAllowMethods, allowedMethods)
			c.SetHeader(fuselage.HeaderAccessControlAllowHeaders, allowedHeaders)
			if config.AllowCredentials {
				c.SetHeader(fuselage.HeaderAccessControlAllowCredentials, "true")
			}

			if c.Request.Method == fuselage.OPTIONS {
				return c.String(http.StatusNoContent, "")
			}

			return next(c)
		}
	}
}
