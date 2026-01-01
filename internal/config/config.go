package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level configuration
type Config struct {
	Scenarios []Scenario `yaml:"scenarios"`
}

// Scenario defines a mock scenario
type Scenario struct {
	Name     string   `yaml:"name"`
	Match    Match    `yaml:"match"`
	Response Response `yaml:"response"`
}

// Match defines the criteria to match a request
type Match struct {
	Path    string            `yaml:"path"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
}

// Response defines the mock response to return
type Response struct {
	Status  int               `yaml:"status"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
	Delay   time.Duration     `yaml:"delay"`
}

// LoadConfig loads the configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
