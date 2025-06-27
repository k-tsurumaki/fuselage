# 🚀 Fuselage

[![CI](https://github.com/k-tsurumaki/fuselage/actions/workflows/ci.yml/badge.svg)](https://github.com/k-tsurumaki/fuselage/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/k-tsurumaki/fuselage)](https://goreportcard.com/report/github.com/k-tsurumaki/fuselage)
[![GoDoc](https://godoc.org/github.com/k-tsurumaki/fuselage?status.svg)](https://godoc.org/github.com/k-tsurumaki/fuselage)

A lightweight, high-performance HTTP web framework for Go, inspired by Echo and Gin but designed with simplicity and developer experience in mind.

## ✨ Features

- **🎯 Simple & Intuitive API** - Echo/Gin-like syntax with improved ergonomics
- **⚡ High Performance** - Zero-allocation routing with custom context management
- **🔧 Built-in Validation** - Struct validation with custom error handling
- **🛡️ Production Ready** - Request ID tracking, panic recovery, timeouts
- **📝 YAML Configuration** - File-based configuration with middleware auto-loading
- **🔀 Route Groups** - Organize routes with prefixes and middleware
- **🚫 Conflict Detection** - Duplicate route registration prevention
- **📊 Request Logging** - Structured logging with request IDs and timing

## 🏆 Why Fuselage?

| Feature | Fuselage | Gin | Echo | Fiber |
|---------|----------|-----|------|-------|
| **Zero Dependencies** | ✅ | ❌ | ❌ | ❌ |
| **Built-in Validation** | ✅ | ❌ | ❌ | ❌ |
| **YAML Config** | ✅ | ❌ | ❌ | ❌ |
| **Route Conflict Detection** | ✅ | ❌ | ❌ | ❌ |
| **Request ID Tracking** | ✅ | Plugin | Plugin | Plugin |
| **Method-specific Routing** | ✅ | ✅ | ✅ | ✅ |
| **Middleware Support** | ✅ | ✅ | ✅ | ✅ |

### 🎯 Fuselage's Unique Strengths

1. **Zero External Dependencies** - Pure Go implementation with no third-party dependencies
2. **Configuration-First Design** - YAML-based configuration with automatic middleware loading
3. **Developer Experience** - Built-in validation, conflict detection, and structured error handling
4. **Production Features** - Request ID tracking, structured logging, and panic recovery out of the box

## 🚀 Quick Start

### Installation

```bash
go get github.com/k-tsurumaki/fuselage
```

### Basic Usage

```go
package main

import (
    "net/http"
    "github.com/k-tsurumaki/fuselage"
)

func main() {
    // Create router
    router := fuselage.New()
    
    // Add middleware
    router.Use(fuselage.RequestID)
    router.Use(fuselage.Logger)
    router.Use(fuselage.Recover)
    
    // Define routes
    router.GET("/", func(c *fuselage.Context) error {
        return c.JSON(http.StatusOK, map[string]string{
            "message": "Hello, Fuselage!",
        })
    })
    
    router.GET("/users/:id", func(c *fuselage.Context) error {
        id, err := c.ParamInt("id")
        if err != nil {
            return c.String(http.StatusBadRequest, "Invalid ID")
        }
        
        return c.JSON(http.StatusOK, map[string]int{
            "user_id": id,
        })
    })
    
    // Start server
    server := fuselage.NewServer(":8080", router)
    server.ListenAndServe()
}
```

## 📋 Configuration-Based Setup

Create `config.yaml`:

```yaml
server:
  host: "localhost"
  port: 8080
  readTimeout: 15s
  writeTimeout: 15s
  idleTimeout: 60s

middleware:
  - requestid
  - logger
  - recover
  - timeout
```

Use configuration:

```go
func main() {
    config, _ := fuselage.LoadConfig("config.yaml")
    router := fuselage.New()
    
    // Routes are automatically configured with middleware
    router.GET("/api/users", getUsers)
    
    server := fuselage.NewServerFromConfig(config, router)
    server.ListenAndServe()
}
```

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   HTTP Request  │───▶│   Fuselage       │───▶│   Your Handler  │
│                 │    │   Router         │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │   Middleware     │
                    │   Chain (LIFO)   │
                    │                  │
                    │ • RequestID      │
                    │ • Logger         │
                    │ • Recover        │
                    │ • Timeout        │
                    │ • Custom...      │
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │   Context        │
                    │                  │
                    │ • Params         │
                    │ • Query          │
                    │ • JSON/String    │
                    │ • Validation     │
                    └──────────────────┘
```

## 🔧 Advanced Features

### Route Groups

```go
router := fuselage.New()

// API v1 group
v1 := router.Group("/api/v1")
v1.GET("/users", listUsers)
v1.POST("/users", createUser)

// Admin group with auth middleware
admin := router.Group("/admin", authMiddleware)
admin.GET("/stats", getStats)
```

### Built-in Validation

```go
type User struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required"`
}

func createUser(c *fuselage.Context) error {
    var user User
    if err := fuselage.Bind(c, &user); err != nil {
        return err // Automatically returns validation errors
    }
    
    // Process valid user...
    return c.JSON(http.StatusCreated, user)
}
```

### Custom Error Handlers

```go
router := fuselage.New()

// Custom 404 handler
router.SetNotFoundHandler(func(c *fuselage.Context) error {
    return c.JSON(http.StatusNotFound, map[string]string{
        "error": "Resource not found",
        "path":  c.Request.URL.Path,
    })
})

// Custom 405 handler
router.SetMethodNotAllowedHandler(func(c *fuselage.Context) error {
    return c.JSON(http.StatusMethodNotAllowed, map[string]string{
        "error": "Method not allowed",
        "method": c.Request.Method,
    })
})
```

### Request Context Methods

```go
func handler(c *fuselage.Context) error {
    // URL parameters
    id := c.Param("id")
    userID, err := c.ParamInt("user_id")
    
    // Query parameters
    page := c.Query("page")
    limit, err := c.QueryInt("limit")
    
    // Headers
    auth := c.GetHeader("Authorization")
    c.Header("X-Custom", "value")
    
    // JSON binding with validation
    var data MyStruct
    if err := fuselage.Bind(c, &data); err != nil {
        return err
    }
    
    // Responses
    return c.JSON(http.StatusOK, data)
    // or
    return c.String(http.StatusOK, "Hello World")
}
```

## 🛠️ Middleware

### Built-in Middleware

- **RequestID** - Adds unique request ID to each request
- **Logger** - Structured request logging with timing
- **Recover** - Panic recovery with logging
- **Timeout** - Request timeout handling

### Custom Middleware

```go
func CustomMiddleware(next fuselage.HandlerFunc) fuselage.HandlerFunc {
    return func(c *fuselage.Context) error {
        // Before request
        start := time.Now()
        
        // Process request
        err := next(c)
        
        // After request
        duration := time.Since(start)
        log.Printf("Request took %v", duration)
        
        return err
    }
}

router.Use(CustomMiddleware)
```

## 📊 Performance

Fuselage is designed for high performance with:

- **Zero-allocation routing** for static routes
- **Efficient parameter extraction** using custom context
- **Minimal memory footprint** with no external dependencies
- **Fast middleware chain** with LIFO execution

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## 📁 Project Structure

```
fuselage/
├── fuselage.go      # Core types and interfaces
├── router.go        # HTTP routing logic
├── context.go       # Request/response context
├── middleware.go    # Built-in middleware
├── validator.go     # Struct validation
├── config.go        # YAML configuration
├── server.go        # HTTP server wrapper
├── types.go         # Custom types
└── example/         # Example applications
    ├── with-config/    # YAML config example
    └── without-config/ # Programmatic example
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [Echo](https://echo.labstack.com/) and [Gin](https://gin-gonic.com/)
- Built with ❤️ for the Go community

---

**Made with 🚀 by the Fuselage team**