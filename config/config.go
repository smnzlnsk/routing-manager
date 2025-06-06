package config

import (
	"fmt"
	"time"
)

// ConfigLoader defines the interface for loading configuration
type ConfigLoader interface {
	Load() (*Config, error)
}

// ConfigLoaderFactory is a factory for creating ConfigLoaders
type ConfigLoaderFactory interface {
	Create(configLoaderType configLoaderType) ConfigLoader
	CreateWithPath(configLoaderType configLoaderType, configPath string) ConfigLoader
}

type configLoaderFactory struct{}

type configLoaderType int

const (
	YamlLoader configLoaderType = iota
	EnvLoader
)

// Config holds the application configuration
type Config struct {
	MonitoringManager MonitoringManagerConfig `yaml:"monitoring_manager"`
	ServiceManager    ServiceManagerConfig    `yaml:"service_manager"`
	MongoDB           MongoDBConfig           `yaml:"mongodb"`
	HTTPServer        HTTPServerConfig        `yaml:"http_server"`
}

type HTTPServerConfig struct {
	Port int `yaml:"port"`
}

// MonitoringManagerConfig holds Monitoring Manager configuration
type MonitoringManagerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// ServiceManagerConfig holds Service Manager configuration
type ServiceManagerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	Host     string        `yaml:"host"`
	Port     int           `yaml:"port"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	Timeout  time.Duration `yaml:"timeout"`
}

type MongoDBDatabaseHandle struct {
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

// NewConfigLoaderFactory creates a new ConfigLoaderFactory
func NewConfigLoaderFactory() ConfigLoaderFactory {
	return &configLoaderFactory{}
}

// Create creates a new ConfigLoader
func (f *configLoaderFactory) Create(configLoaderType configLoaderType) ConfigLoader {
	switch configLoaderType {
	case YamlLoader:
		return NewYAMLConfigLoader("config.yaml")
	case EnvLoader:
		return NewEnvConfigLoader(".env")
	default:
		return nil
	}
}

// CreateWithPath creates a new ConfigLoader with a specific path
func (f *configLoaderFactory) CreateWithPath(configLoaderType configLoaderType, configPath string) ConfigLoader {
	switch configLoaderType {
	case YamlLoader:
		return NewYAMLConfigLoader(configPath)
	case EnvLoader:
		return NewEnvConfigLoader(configPath)
	default:
		return nil
	}
}

// validateConfig validates the configuration
func validateConfig(cfg *Config) error {
	if cfg.MonitoringManager.Host == "" {
		return fmt.Errorf("monitoring manager host is required")
	}

	return nil
}

// setDefaults sets default values for empty configuration fields
func setDefaults(cfg *Config) {
	// MonitoringManager defaults
	if cfg.MonitoringManager.Host == "" {
		cfg.MonitoringManager.Host = "monitoring_manager"
	}
	if cfg.MonitoringManager.Port == 0 {
		cfg.MonitoringManager.Port = 10999
	}

	// ServiceManager defaults
	if cfg.ServiceManager.Host == "" {
		cfg.ServiceManager.Host = "cluster_service_manager"
	}
	if cfg.ServiceManager.Port == 0 {
		cfg.ServiceManager.Port = 10110
	}

	// MongoDB defaults
	if cfg.MongoDB.Host == "" {
		cfg.MongoDB.Host = "mongo_clusternet"
	}
	if cfg.MongoDB.Port == 0 {
		cfg.MongoDB.Port = 10108
	}
	if cfg.MongoDB.Username == "" {
		cfg.MongoDB.Username = ""
	}
	if cfg.MongoDB.Password == "" {
		cfg.MongoDB.Password = ""
	}
	if cfg.MongoDB.Timeout == 0 {
		cfg.MongoDB.Timeout = 30 * time.Second
	}
}
