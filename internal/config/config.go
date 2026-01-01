package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Scenarios []Scenario `yaml:"scenarios"`
}

type Scenario struct {
	Name     string   `yaml:"name"`
	Match    Match    `yaml:"match"`
	Response Response `yaml:"response"`
}

type Match struct {
	Path    string            `yaml:"path"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
}

type Response struct {
	Status  int               `yaml:"status"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
	Delay   time.Duration     `yaml:"delay"`
}

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
