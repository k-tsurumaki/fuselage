package fuselage

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Router handles HTTP routing with middleware support
type Router struct {
	routes                  map[string]map[string]HandlerFunc
	middleware              []MiddlewareFunc
	notFoundHandler         HandlerFunc
	methodNotAllowedHandler HandlerFunc
	prefix                  string
}

// New creates a new Router instance
func New() *Router {
	return &Router{
		routes:                  make(map[string]map[string]HandlerFunc),
		notFoundHandler:         defaultNotFound,
		methodNotAllowedHandler: defaultMethodNotAllowed,
	}
}

// Use adds middleware to the router (LIFO order)
func (r *Router) Use(middleware MiddlewareFunc) {
	r.middleware = append(r.middleware, middleware)
}

// Group creates a route group with prefix and middleware
func (r *Router) Group(prefix string, middlewares ...MiddlewareFunc) *Router {
	group := &Router{
		routes:                  make(map[string]map[string]HandlerFunc),
		middleware:              append(r.middleware, middlewares...),
		notFoundHandler:         r.notFoundHandler,
		methodNotAllowedHandler: r.methodNotAllowedHandler,
		prefix:                  r.prefix + prefix,
	}
	return group
}

// SetNotFoundHandler sets custom 404 handler
func (r *Router) SetNotFoundHandler(handler HandlerFunc) {
	r.notFoundHandler = handler
}

// SetMethodNotAllowedHandler sets custom 405 handler
func (r *Router) SetMethodNotAllowedHandler(handler HandlerFunc) {
	r.methodNotAllowedHandler = handler
}

// ServeHTTP implements http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Request:  req,
		Response: w,
		status:   http.StatusOK,
	}

	handler, params := r.findHandler(req.Method, req.URL.Path)
	if handler == nil {
		if r.hasPath(req.URL.Path) {
			handler = r.methodNotAllowedHandler
		} else {
			handler = r.notFoundHandler
		}
	}

	ctx.params = params
	finalHandler := r.applyMiddleware(handler)
	
	if err := finalHandler(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *Router) applyMiddleware(handler HandlerFunc) HandlerFunc {
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	return handler
}

func (r *Router) addRoute(method, path string, handler HandlerFunc) error {
	if !isValidPath(path) {
		return errors.New("invalid path")
	}

	fullPath := r.prefix + path
	key := method + " " + fullPath

	if r.routes[method] == nil {
		r.routes[method] = make(map[string]HandlerFunc)
	}

	if _, exists := r.routes[method][fullPath]; exists {
		return fmt.Errorf("route %s already exists", key)
	}

	r.routes[method][fullPath] = handler
	return nil
}

func (r *Router) findHandler(method, path string) (handler HandlerFunc, params map[string]string) {
	if methodRoutes, exists := r.routes[method]; exists {
		if h, found := methodRoutes[path]; found {
			return h, nil
		}

		for routePath, h := range methodRoutes {
			if p := matchRoute(routePath, path); p != nil {
				return h, p
			}
		}
	}
	return nil, nil
}

func (r *Router) hasPath(path string) bool {
	for _, methodRoutes := range r.routes {
		for routePath := range methodRoutes {
			if routePath == path || matchRoute(routePath, path) != nil {
				return true
			}
		}
	}
	return false
}

// GET registers a GET route
func (r *Router) GET(path string, handler HandlerFunc) error {
	return r.addRoute("GET", path, handler)
}

// POST registers a POST route
func (r *Router) POST(path string, handler HandlerFunc) error {
	return r.addRoute("POST", path, handler)
}

// PUT registers a PUT route
func (r *Router) PUT(path string, handler HandlerFunc) error {
	return r.addRoute("PUT", path, handler)
}

// DELETE registers a DELETE route
func (r *Router) DELETE(path string, handler HandlerFunc) error {
	return r.addRoute("DELETE", path, handler)
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

func defaultNotFound(c *Context) error {
	return c.String(http.StatusNotFound, "Not Found")
}

func defaultMethodNotAllowed(c *Context) error {
	return c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
}

func isValidPath(path string) bool {
	return path != "" && strings.HasPrefix(path, "/")
}