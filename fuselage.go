package fuselage

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// Router wraps http.ServeMux with middleware support
type Router struct {
	mux        *http.ServeMux
	middleware []func(http.Handler) http.Handler
	paramRoutes map[string]paramRoute
}

type paramRoute struct {
	handler http.HandlerFunc
	pattern string
}

// New creates a new Router instance
func New() *Router {
	return &Router{
		mux:         http.NewServeMux(),
		paramRoutes: make(map[string]paramRoute),
	}
}

// Use adds middleware to the router
func (r *Router) Use(middleware func(http.Handler) http.Handler) {
	r.middleware = append(r.middleware, middleware)
}

// ServeHTTP implements http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Check for parameterized routes first
	for key, route := range r.paramRoutes {
		if strings.HasPrefix(key, req.Method+" ") {
			if params := matchRoute(route.pattern, req.URL.Path); params != nil {
				handler := &paramHandler{handler: route.handler, params: params}
				r.applyMiddleware(handler).ServeHTTP(w, req)
				return
			}
		}
	}

	// Fall back to ServeMux
	r.applyMiddleware(r.mux).ServeHTTP(w, req)
}

func (r *Router) applyMiddleware(handler http.Handler) http.Handler {
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	return handler
}

// GET registers a GET route
func (r *Router) GET(path string, handler http.HandlerFunc) {
	if strings.Contains(path, ":") {
		r.paramRoutes["GET "+path] = paramRoute{handler: handler, pattern: path}
	} else {
		r.mux.HandleFunc("GET "+path, handler)
	}
}

// POST registers a POST route
func (r *Router) POST(path string, handler http.HandlerFunc) {
	if strings.Contains(path, ":") {
		r.paramRoutes["POST "+path] = paramRoute{handler: handler, pattern: path}
	} else {
		r.mux.HandleFunc("POST "+path, handler)
	}
}

// PUT registers a PUT route
func (r *Router) PUT(path string, handler http.HandlerFunc) {
	if strings.Contains(path, ":") {
		r.paramRoutes["PUT "+path] = paramRoute{handler: handler, pattern: path}
	} else {
		r.mux.HandleFunc("PUT "+path, handler)
	}
}

// DELETE registers a DELETE route
func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	if strings.Contains(path, ":") {
		r.paramRoutes["DELETE "+path] = paramRoute{handler: handler, pattern: path}
	} else {
		r.mux.HandleFunc("DELETE "+path, handler)
	}
}

func matchRoute(routePath, requestPath string) map[string]string {
	routeParts := strings.Split(routePath, "/")
	requestParts := strings.Split(requestPath, "/")

	if len(routeParts) != len(requestParts) {
		return nil
	}

	params := make(map[string]string)
	for i, part := range routeParts {
		if strings.HasPrefix(part, ":") {
			params[part[1:]] = requestParts[i]
		} else if part != requestParts[i] {
			return nil
		}
	}

	return params
}

type paramHandler struct {
	handler http.HandlerFunc
	params  map[string]string
}

func (ph *paramHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	for k, v := range ph.params {
		ctx = context.WithValue(ctx, k, v)
	}
	ph.handler(w, r.WithContext(ctx))
}

// GetParam extracts URL parameter from request context
func GetParam(r *http.Request, key string) string {
	if val := r.Context().Value(key); val != nil {
		return val.(string)
	}
	return ""
}

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