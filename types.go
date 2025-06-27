package fuselage

// ParamKey is a custom type for context keys to avoid collisions
type ParamKey string

// RequestIDKey is the context key for request ID
const RequestIDKey ParamKey = "request_id"