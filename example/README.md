# 🚀 Fuselage Example

This directory contains a complete example application demonstrating how to use Fuselage.

## 🚀 Quick Start

Run the example application:

```bash
go run main.go
```

**Features:**
- Port: 8082
- Middleware: RequestID, Logger, Recover, Timeout
- Built-in validation with error handling
- Complete REST API implementation

## 🌐 API Endpoints

The example provides a complete REST API:

| Method | Endpoint | Description | Validation |
|--------|----------|-------------|------------|
| `GET` | `/users` | List all users | - |
| `GET` | `/users/:id` | Get user by ID | ID must be integer |
| `POST` | `/users` | Create new user | Name required, min 2 chars |
| `PUT` | `/users/:id` | Update user | ID + Name validation |
| `DELETE` | `/users/:id` | Delete user | ID must be integer |

## 🧪 Testing the API

```bash
# Get all users
curl http://localhost:8082/users

# Get specific user
curl http://localhost:8082/users/1

# Create new user
curl -X POST http://localhost:8082/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Charlie"}'

# Update user
curl -X PUT http://localhost:8082/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated"}'

# Delete user
curl -X DELETE http://localhost:8082/users/2
```