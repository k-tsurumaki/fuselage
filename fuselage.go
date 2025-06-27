package fuselage

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// Router handles HTTP routing
type Router struct {
	routes     map[string]map[string]http.HandlerFunc
	middleware []func(http.Handler) http.Handler
}

// New creates a new Router instance
func New() *Router {
	return &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}

// Use adds middleware to the router
func (r *Router) Use(middleware func(http.Handler) http.Handler) {
	r.middleware = append(r.middleware, middleware)
}

// ServeHTTP implements http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := r.findHandler(req.Method, req.URL.Path)
	if handler == nil {
		http.NotFound(w, req)
		return
	}

	// Apply middleware in reverse order
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}

	handler.ServeHTTP(w, req)
}

// GET registers a GET route
func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.addRoute("GET", path, handler)
}

// POST registers a POST route
func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.addRoute("POST", path, handler)
}

// PUT registers a PUT route
func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.addRoute("PUT", path, handler)
}

// DELETE registers a DELETE route
func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.addRoute("DELETE", path, handler)
}

func (r *Router) addRoute(method, path string, handler http.HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]http.HandlerFunc)
	}
	r.routes[method][path] = handler
}

func (r *Router) findHandler(method, path string) http.Handler {
	if methodRoutes, exists := r.routes[method]; exists {
		if handler, found := methodRoutes[path]; found {
			return handler
		}

		// Check for parameterized routes
		for routePath, handler := range methodRoutes {
			if params := matchRoute(routePath, path); params != nil {
				return &paramHandler{handler: handler, params: params}
			}
		}
	}
	return nil
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