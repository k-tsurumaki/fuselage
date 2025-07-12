# 🚀 Fuselage

[![CI](https://github.com/k-tsurumaki/fuselage/actions/workflows/ci.yml/badge.svg)](https://github.com/k-tsurumaki/fuselage/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/k-tsurumaki/fuselage)](https://goreportcard.com/report/github.com/k-tsurumaki/fuselage)
[![GoDoc](https://godoc.org/github.com/k-tsurumaki/fuselage?status.svg)](https://godoc.org/github.com/k-tsurumaki/fuselage)

A lightweight, high-performance HTTP web framework for Go, inspired by Echo and Gin but designed with simplicity and developer experience in mind.

**Current Version: v1.0.0** - Production ready with stable API

## ✨ Features

- **🎯 Simple & Intuitive API** - Echo/Gin-like syntax with improved ergonomics
- **⚡ High Performance** - Zero-allocation routing with custom context management
- **🔧 Built-in Validation** - Struct validation with custom error handling
- **🛡️ Production Ready** - Request ID tracking, panic recovery, timeouts

- **🔀 Route Groups** - Organize routes with prefixes and middleware
- **🚫 Conflict Detection** - Duplicate route registration prevention
- **📊 Request Logging** - Structured logging with request IDs and timing
- **🧩 Modular Middleware** - Separate configurable middleware package
- **🌐 Advanced CORS** - Pattern matching and credential support
- **📜 Rich Constants** - Comprehensive HTTP methods and headers

## 🏆 Why Fuselage?

| Feature | Fuselage | Gin | Echo | Fiber |
|---------|----------|-----|------|-------|
| **Zero Dependencies** | ✅ | ❌ | ❌ | ❌ |
| **Built-in Validation** | ✅ | ❌ | ❌ | ❌ |
| **Route Conflict Detection** | ✅ | ❌ | ❌ | ❌ |
| **Request ID Tracking** | ✅ | Plugin | Plugin | Plugin |
| **Configurable Middleware** | ✅ | ❌ | ✅ | ✅ |
| **CORS with Patterns** | ✅ | Plugin | Plugin | Plugin |
| **HTTP Constants** | ✅ | ❌ | ✅ | ✅ |
| **Method-specific Routing** | ✅ | ✅ | ✅ | ✅ |
| **Middleware Support** | ✅ | ✅ | ✅ | ✅ |

### 🎯 Fuselage's Unique Strengths

1. **Zero External Dependencies** - Pure Go implementation with no third-party dependencies
2. **Developer Experience** - Built-in validation, conflict detection, and structured error handling
4. **Production Features** - Request ID tracking, structured logging, and panic recovery out of the box
5. **Modular Middleware** - Separate configurable middleware package with advanced features
6. **Rich Constants** - Comprehensive HTTP method and header constants for better code clarity
7. **Enhanced Context** - Improved context methods with better status and header management

## 🚀 Quick Start

### Installation

```bash
go get github.com/k-tsurumaki/fuselage@v1.0.0
```

### Basic Usage

```go
package main

import (
    "net/http"
    "github.com/k-tsurumaki/fuselage"
    "github.com/k-tsurumaki/fuselage/middleware"
)

func main() {
    // Create router
    router := fuselage.New()
    
    // Add middleware from middleware package
    router.Use(middleware.RequestID())
    router.Use(middleware.Logger())
    router.Use(middleware.Recover())
    router.Use(middleware.CORS())
    
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

## ⚙️ Server Configuration

```go
func main() {
    router := fuselage.New()
    
    // Add middleware
    router.Use(middleware.RequestID())
    router.Use(middleware.Logger())
    router.Use(middleware.Recover())
    router.Use(middleware.Timeout())
    
    // Define routes
    router.GET("/api/users", getUsers)
    
    // Create server with custom settings
    server := fuselage.NewServer(":8080", router)
    server.ReadTimeout = 15 * time.Second
    server.WriteTimeout = 15 * time.Second
    server.IdleTimeout = 60 * time.Second
    
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

### Enhanced Context Methods

```go
func handler(c *fuselage.Context) error {
    // URL parameters
    id := c.Param("id")
    userID, err := c.ParamInt("user_id")
    
    // Query parameters
    page := c.Query("page")
    limit, err := c.QueryInt("limit")
    
    // Headers (updated methods)
    auth := c.Header("Authorization")        // Get request header
    c.SetHeader("X-Custom", "value")        // Set response header
    
    // Status handling
    c.SetStatus(http.StatusCreated)          // Set status
    currentStatus := c.Status()              // Get current status
    
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

### Built-in Middleware Package

Fuselage now includes a separate `middleware` package with configurable middleware:

```go
import "github.com/k-tsurumaki/fuselage/middleware"

// Basic usage
router.Use(middleware.RequestID())
router.Use(middleware.Logger())
router.Use(middleware.Recover())
router.Use(middleware.Timeout())
router.Use(middleware.CORS())
```

### Configurable Middleware

```go
// Custom RequestID with generator
router.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
    Generator: func() string {
        return "custom-" + uuid.New().String()
    },
    TargetHeader: "X-Custom-Request-ID",
}))

// Custom CORS configuration
router.Use(middleware.CORSWithConfig(&middleware.CORSConfig{
    AllowedOrigins:   []string{"https://example.com", "https://*.example.com"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders:   []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
}))

// Custom timeout
router.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
    Timeout: 60 * time.Second,
}))
```

### Available Middleware

- **RequestID** - Unique request ID generation and tracking
- **Logger** - Structured access logging with request IDs
- **Recover** - Panic recovery with detailed logging
- **Timeout** - Configurable request timeout handling
- **CORS** - Cross-Origin Resource Sharing with pattern matching

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
├── fuselage.go         # Core types, HTTP methods, and header constants
├── router.go           # HTTP routing logic
├── context.go          # Enhanced request/response context
├── validator.go        # Struct validation

├── server.go           # HTTP server wrapper
├── types.go            # Custom types and constants
├── middleware/         # Middleware package
│   ├── accessLogger.go # Access logging middleware
│   ├── cors.go         # CORS middleware with pattern matching
│   ├── recover.go      # Panic recovery middleware
│   ├── requestID.go    # Request ID generation and tracking
│   └── timeout.go      # Request timeout middleware
├── templates/          # Code generation templates
│   ├── adapter/        # Adapter pattern templates
│   ├── domain/         # Domain layer templates
│   └── service/        # Service layer templates
└── example/            # Example application
    ├── main.go         # Complete REST API example
    ├── go.mod          # Module definition
    └── README.md       # Usage instructions
```

## 📋 Versioning

Fuselage follows [Semantic Versioning](https://semver.org/). For the versions available, see the [tags on this repository](https://github.com/k-tsurumaki/fuselage/tags).

### Version History
- **v1.0.0** - Initial stable release with full feature set

## 🚀 Releasing

For maintainers releasing new versions, see [RELEASE.md](RELEASE.md) for detailed release procedures.

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