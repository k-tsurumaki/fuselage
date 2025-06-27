package fuselage

import (
	"net/http"
	"time"
)

// Server wraps http.Server
type Server struct {
	*http.Server
}

// NewServer creates a new Server instance
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}