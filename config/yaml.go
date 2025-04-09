package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// YAMLConfigLoader loads configuration from a YAML file
type YAMLConfigLoader struct {
	filePath string
}

// NewYAMLConfigLoader creates a new YAML config loader
func NewYAMLConfigLoader(filePath string) *YAMLConfigLoader {
	return &YAMLConfigLoader{
		filePath: filePath,
	}
}

// Load loads the configuration from a YAML file
func (y *YAMLConfigLoader) Load() (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(y.filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", y.filePath)
	}

	// Read file
	data, err := os.ReadFile(y.filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing YAML: %w", err)
	}

	// Set defaults for empty values
	setDefaults(&cfg)

	// Validate configuration
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
