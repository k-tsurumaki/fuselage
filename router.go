package fuselage

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Router handles HTTP routing with middleware support
type Router struct {
	routes                  map[string]map[string]routeEntry
	middleware              []MiddlewareFunc
	notFoundHandler         HandlerFunc
	methodNotAllowedHandler HandlerFunc
	prefix                  string
}

type routeEntry struct {
	handler     HandlerFunc
	middlewares []MiddlewareFunc
}

// New creates a new Router instance
func New() *Router {
	return &Router{
		routes:                  make(map[string]map[string]routeEntry),
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
		routes:                  make(map[string]map[string]routeEntry),
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
		status:   0,
		written:  false,
	}

	handler, params, routeMiddlewares := r.findHandler(req.Method, req.URL.Path)
	if handler == nil {
		if r.hasPath(req.URL.Path) {
			handler = r.methodNotAllowedHandler
		} else {
			handler = r.notFoundHandler
		}
	}

	ctx.params = params
	finalHandler := r.applyMiddlewareWithRoute(handler, routeMiddlewares)

	if err := finalHandler(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *Router) applyMiddlewareWithRoute(handler HandlerFunc, routeMiddlewares []MiddlewareFunc) HandlerFunc {
	// global → group → route
	all := append([]MiddlewareFunc{}, r.middleware...)
	all = append(all, routeMiddlewares...)
	for i := len(all) - 1; i >= 0; i-- {
		handler = all[i](handler)
	}
	return handler
}

func (r *Router) addRoute(method, path string, handler HandlerFunc, middlewares ...MiddlewareFunc) error {
	if !isValidPath(path) {
		return errors.New("invalid path")
	}

	fullPath := r.prefix + path
	key := method + " " + fullPath

	if r.routes[method] == nil {
		r.routes[method] = make(map[string]routeEntry)
	}

	if _, exists := r.routes[method][fullPath]; exists {
		return fmt.Errorf("route %s already exists", key)
	}

	r.routes[method][fullPath] = routeEntry{
		handler:     handler,
		middlewares: middlewares,
	}
	return nil
}

// GET registers a GET route
func (r *Router) GET(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) error {
	return r.addRoute(GET, path, handler, middlewares...)
}

// POST registers a POST route
func (r *Router) POST(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) error {
	return r.addRoute(POST, path, handler, middlewares...)
}

// PUT registers a PUT route
func (r *Router) PUT(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) error {
	return r.addRoute(PUT, path, handler, middlewares...)
}

// DELETE registers a DELETE route
func (r *Router) DELETE(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) error {
	return r.addRoute(DELETE, path, handler, middlewares...)
}

// findHandler locates the handler, path parameters, and route-specific middlewares for a given HTTP method and path.
func (r *Router) findHandler(method, path string) (HandlerFunc, map[string]string, []MiddlewareFunc) {
	if methodRoutes, exists := r.routes[method]; exists {
		// If there is an exact match, return it
		if entry, found := methodRoutes[path]; found {
			return entry.handler, nil, entry.middlewares
		}
		// Otherwise, try to match with path parameters (e.g., /users/:id)
		for routePath, entry := range methodRoutes {
			if p := matchRoute(routePath, path); p != nil {
				return entry.handler, p, entry.middlewares
			}
		}
	}
	// No match found
	return nil, nil, nil
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
