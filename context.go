package fuselage

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// Context provides request/response handling
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	params   map[string]string
	status   int
}

// Param gets URL parameter
func (c *Context) Param(key string) string {
	if c.params == nil {
		return ""
	}
	return c.params[key]
}

// ParamInt gets URL parameter as integer
func (c *Context) ParamInt(key string) (value int, err error) {
	str := c.Param(key)
	if str == "" {
		return 0, errors.New("parameter not found")
	}
	return strconv.Atoi(str)
}

// Query gets query parameter
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// QueryInt gets query parameter as integer
func (c *Context) QueryInt(key string) (value int, err error) {
	str := c.Query(key)
	if str == "" {
		return 0, errors.New("query parameter not found")
	}
	return strconv.Atoi(str)
}

// Bind binds JSON request body to struct
func (c *Context) Bind(v interface{}) error {
	if c.Request.Body == nil {
		return errors.New("request body is empty")
	}
	return json.NewDecoder(c.Request.Body).Decode(v)
}

// JSON sends JSON response
func (c *Context) JSON(status int, data interface{}) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(status)
	c.status = status
	return json.NewEncoder(c.Response).Encode(data)
}

// String sends string response
func (c *Context) String(status int, text string) error {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(status)
	c.status = status
	_, err := c.Response.Write([]byte(text))
	return err
}

// SetStatus sets response status
func (c *Context) SetStatus(status int) {
	c.Response.WriteHeader(status)
	c.status = status
}

// Status gets response status
func (c *Context) Status() int {
	return c.status
}

// Header sets response header
func (c *Context) SetHeader(key, value string) {
	c.Response.Header().Set(key, value)
}

// Header gets request header
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// Bind validates and binds JSON request body
func Bind(c *Context, v interface{}) error {
	if err := c.Bind(v); err != nil {
		return err
	}

	if errors := ValidateStruct(v); len(errors) > 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"errors": errors,
		})
	}

	return nil
}
