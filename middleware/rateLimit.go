package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/k-tsurumaki/fuselage"
)

type RateLimitConfig struct {
	// Requests per window
	Limit int
	// Time window duration
	Window time.Duration
	// Key generator function (default: IP-based)
	KeyGenerator func(*fuselage.Context) string
	// Skip function to bypass rate limiting
	Skipper func(*fuselage.Context) bool
	// Error handler for rate limit exceeded
	ErrorHandler func(*fuselage.Context) error
}

type rateLimiter struct {
	requests map[string]*bucket
	mutex    sync.RWMutex
	config   RateLimitConfig
}

type bucket struct {
	count     int
	resetTime time.Time
}

var DefaultRateLimitConfig = RateLimitConfig{
	Limit:  100,
	Window: time.Minute,
	KeyGenerator: func(c *fuselage.Context) string {
		ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
		return ip
	},
	Skipper: func(c *fuselage.Context) bool {
		return false
	},
	ErrorHandler: func(c *fuselage.Context) error {
		return c.JSON(http.StatusTooManyRequests, map[string]string{
			"error": "Rate limit exceeded",
		})
	},
}

func RateLimit() fuselage.MiddlewareFunc {
	return RateLimitWithConfig(DefaultRateLimitConfig)
}

func RateLimitWithConfig(config RateLimitConfig) fuselage.MiddlewareFunc {
	if config.Limit <= 0 {
		config.Limit = DefaultRateLimitConfig.Limit
	}
	if config.Window <= 0 {
		config.Window = DefaultRateLimitConfig.Window
	}
	if config.KeyGenerator == nil {
		config.KeyGenerator = DefaultRateLimitConfig.KeyGenerator
	}
	if config.Skipper == nil {
		config.Skipper = DefaultRateLimitConfig.Skipper
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = DefaultRateLimitConfig.ErrorHandler
	}

	limiter := &rateLimiter{
		requests: make(map[string]*bucket),
		config:   config,
	}

	// Cleanup goroutine
	go limiter.cleanup()

	return func(next fuselage.HandlerFunc) fuselage.HandlerFunc {
		return func(c *fuselage.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			key := config.KeyGenerator(c)
			if !limiter.allow(key) {
				return config.ErrorHandler(c)
			}

			return next(c)
		}
	}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	b, exists := rl.requests[key]

	if !exists || now.After(b.resetTime) {
		rl.requests[key] = &bucket{
			count:     1,
			resetTime: now.Add(rl.config.Window),
		}
		return true
	}

	if b.count >= rl.config.Limit {
		return false
	}

	b.count++
	return true
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.Window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for key, bucket := range rl.requests {
			if now.After(bucket.resetTime) {
				delete(rl.requests, key)
			}
		}
		rl.mutex.Unlock()
	}
}