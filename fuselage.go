package fuselage

// HandlerFunc defines the handler function signature
type HandlerFunc func(*Context) error

// MiddlewareFunc defines middleware function signature
type MiddlewareFunc func(HandlerFunc) HandlerFunc