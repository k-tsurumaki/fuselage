package fuselage

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Create temporary config file
	configContent := `server:
  host: "0.0.0.0"
  port: 9000
  readTimeout: 30s
  writeTimeout: 30s
  idleTimeout: 120s

middleware:
  - logger
  - recover`

	tmpFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.Server.Host != "0.0.0.0" {
		t.Errorf("Expected host '0.0.0.0', got '%s'", config.Server.Host)
	}

	if config.Server.Port != 9000 {
		t.Errorf("Expected port 9000, got %d", config.Server.Port)
	}

	if config.Server.ReadTimeout != 30*time.Second {
		t.Errorf("Expected readTimeout 30s, got %v", config.Server.ReadTimeout)
	}

	if len(config.Middleware) != 2 {
		t.Errorf("Expected 2 middleware, got %d", len(config.Middleware))
	}
}

func TestConfigDefaults(t *testing.T) {
	// Create empty config file
	tmpFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString("{}"); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.Server.Host != "localhost" {
		t.Errorf("Expected default host 'localhost', got '%s'", config.Server.Host)
	}

	if config.Server.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", config.Server.Port)
	}
}
