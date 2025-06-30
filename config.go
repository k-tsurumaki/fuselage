package fuselage

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host         string        `yaml:"host"`
		Port         int           `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"readTimeout"`
		WriteTimeout time.Duration `yaml:"writeTimeout"`
		IdleTimeout  time.Duration `yaml:"idleTimeout"`
	} `yaml:"server"`
	Middleware []string `yaml:"middleware"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set defaults
	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 15 * time.Second
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 15 * time.Second
	}
	if config.Server.IdleTimeout == 0 {
		config.Server.IdleTimeout = 60 * time.Second
	}

	return &config, nil
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
