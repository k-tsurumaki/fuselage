# 🚀 Fuselage Examples

This directory contains example applications demonstrating different ways to use Fuselage.

## 📁 Examples

### 🔧 with-config

Demonstrates YAML configuration-based setup with automatic middleware loading.

```bash
cd with-config
go run main.go
```

**Features:**
- Port: 8083
- Configuration: `config.yaml`
- Middleware: Auto-loaded from config
- Built-in validation with error handling

### ⚡ without-config

Demonstrates programmatic setup without configuration files.

```bash
cd without-config
go run main.go
```

**Features:**
- Port: 8082
- Configuration: Programmatic
- Middleware: Manually configured
- Direct middleware chain setup

## 🌐 API Endpoints

Both examples provide the same REST API:

| Method | Endpoint | Description | Validation |
|--------|----------|-------------|------------|
| `GET` | `/users` | List all users | - |
| `GET` | `/users/:id` | Get user by ID | ID must be integer |
| `POST` | `/users` | Create new user | Name required, min 2 chars |
| `PUT` | `/users/:id` | Update user | ID + Name validation |
| `DELETE` | `/users/:id` | Delete user | ID must be integer |

## 🧪 Testing the APIs

```bash
# Get all users
curl http://localhost:8083/users

# Get specific user
curl http://localhost:8083/users/1

# Create new user
curl -X POST http://localhost:8083/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Charlie"}'

# Update user
curl -X PUT http://localhost:8083/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated"}'

# Delete user
curl -X DELETE http://localhost:8083/users/2
```

## 🔍 Key Differences

| Feature | with-config | without-config |
|---------|-------------|----------------|
| **Setup** | YAML-driven | Code-driven |
| **Middleware** | Auto-loaded | Manual setup |
| **Configuration** | External file | In-code |
| **Flexibility** | Template-based | Full control |
| **Best for** | Production | Development |