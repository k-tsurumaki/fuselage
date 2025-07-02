package middleware

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/k-tsurumaki/fuselage"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int // seconds
}

var DefaultCORSConfig = CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{fuselage.GET, fuselage.HEAD, fuselage.PUT, fuselage.PATCH, fuselage.POST, fuselage.DELETE},
	MaxAge:       600, // 10分
}

func CORS() fuselage.MiddlewareFunc {
	return CORSWithConfig(&DefaultCORSConfig)
}

func CORSWithConfig(config *CORSConfig) fuselage.MiddlewareFunc {
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = DefaultCORSConfig.AllowOrigins
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = DefaultCORSConfig.AllowMethods
	}
	if len(config.AllowHeaders) == 0 {
		config.AllowHeaders = []string{fuselage.HeaderContentType, fuselage.HeaderAuthorization}
	}

	allowOriginPatterns := make([]*regexp.Regexp, 0, len(config.AllowOrigins))
	for _, origin := range config.AllowOrigins {
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
		allowOriginPatterns = append(allowOriginPatterns, re)
	}

	allowMethods := strings.Join(config.AllowMethods, ", ")
	allowHeaders := strings.Join(config.AllowHeaders, ", ")
	exposeHeaders := strings.Join(config.ExposeHeaders, ", ")
	maxAge := ""
	if config.MaxAge > 0 {
		maxAge = strconv.Itoa(config.MaxAge)
	}

	return func(next fuselage.HandlerFunc) fuselage.HandlerFunc {
		return func(c *fuselage.Context) error {
			origin := c.Request.Header.Get(fuselage.HeaderOrigin)
			allow := false

			c.SetHeader(fuselage.HeaderVary, fuselage.HeaderOrigin)

			preflight := c.Request.Method == fuselage.OPTIONS

			if origin != "" {
				for _, o := range config.AllowOrigins {
					if o == "*" || o == origin {
						allow = true
						break
					}
				}
				if !allow {
					for _, re := range allowOriginPatterns {
						if re.MatchString(origin) {
							allow = true
							break
						}
					}
				}
				if allow {
					c.SetHeader(fuselage.HeaderAccessControlAllowOrigin, origin)
				}
			}

			c.SetHeader(fuselage.HeaderAccessControlAllowMethods, allowMethods)
			c.SetHeader(fuselage.HeaderAccessControlAllowHeaders, allowHeaders)
			if exposeHeaders != "" {
				c.SetHeader(fuselage.HeaderAccessControlExposeHeaders, exposeHeaders)
			}
			if maxAge != "" {
				c.SetHeader(fuselage.HeaderAccessControlMaxAge, maxAge)
			}
			if config.AllowCredentials {
				c.SetHeader(fuselage.HeaderAccessControlAllowCredentials, "true")
			}

			if preflight {
				// Verify Access-Control-Request-Method/Headers
				reqMethod := c.Request.Header.Get(fuselage.HeaderAccessControlRequestMethod)
				if reqMethod != "" && !contains(config.AllowMethods, reqMethod) {
					return c.String(http.StatusForbidden, "CORS: method not allowed")
				}
				reqHeaders := c.Request.Header.Get(fuselage.HeaderAccessControlRequestHeaders)
				if reqHeaders != "" {
					reqHeaderList := strings.Split(reqHeaders, ",")
					for _, h := range reqHeaderList {
						h = strings.TrimSpace(h)
						if h != "" && !contains(config.AllowHeaders, h) {
							return c.String(http.StatusForbidden, "CORS: header not allowed")
						}
					}
				}
				return c.String(http.StatusNoContent, "")
			}

			return next(c)
		}
	}
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if strings.EqualFold(s, v) {
			return true
		}
	}
	return false
}
