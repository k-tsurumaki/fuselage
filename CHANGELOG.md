# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.0.0] - 2025-06-30

### Added
- 🎯 Simple & Intuitive API with Echo/Gin-like syntax
- ⚡ High Performance zero-allocation routing
- 🔧 Built-in struct validation with custom error handling
- 🛡️ Production-ready features (Request ID, panic recovery, timeouts)
- 🔀 Route groups with prefix and middleware support
- 🚫 Duplicate route registration detection
- 📊 Structured request logging with timing
- 🧩 Modular middleware package with configurable options
- 🌐 Advanced CORS with pattern matching and credentials
- 📜 Comprehensive HTTP method and header constants
- 🎨 Enhanced Context with improved status and header management

### Features
- Zero external dependencies
- Method-specific routing (GET, POST, PUT, DELETE, etc.)
- URL parameter extraction with type conversion
- Query parameter handling
- JSON request/response binding
- Custom error handlers (404, 405)
- Middleware chain with LIFO execution
- Request timeout handling
- Panic recovery with logging
- Request ID generation and tracking
- CORS support with origin pattern matching

### Middleware Package
- `RequestID()` - Unique request ID generation
- `Logger()` - Access logging with request IDs
- `Recover()` - Panic recovery with detailed logging
- `Timeout()` - Configurable request timeouts
- `CORS()` - Cross-origin resource sharing

### Initial Release
This is the first stable release of Fuselage, ready for production use.
