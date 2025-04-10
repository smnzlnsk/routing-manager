package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// EnvConfigLoader loads configuration from environment variables
type EnvConfigLoader struct {
	dotEnvPath string
}

// NewEnvConfigLoader creates a new environment variable config loader
func NewEnvConfigLoader(dotEnvPath string) *EnvConfigLoader {
	return &EnvConfigLoader{
		dotEnvPath: dotEnvPath,
	}
}

// Load loads the configuration from environment variables
func (e *EnvConfigLoader) Load() (*Config, error) {
	// Load .env file if it exists and path is provided
	if e.dotEnvPath != "" {
		_ = godotenv.Load(e.dotEnvPath)
	} else {
		_ = godotenv.Load() // Try default .env in current directory
	}

	cfg := &Config{
		HTTPServer: HTTPServerConfig{
			Port: getEnvAsInt("ROUTING_MANAGER_HTTP_SERVER_PORT", 8080),
		},
		MQTT: MQTTConfig{
			Host:           getEnv("MQTT_BROKER_HOST", "mqtt"),
			Port:           getEnvAsInt("MQTT_BROKER_PORT", 1883),
			ClientID:       getEnv("MQTT_CLIENT_ID", "mqtt-worker-"+strconv.FormatInt(time.Now().UnixNano(), 10)),
			Username:       getEnv("MQTT_USERNAME", ""),
			Password:       getEnv("MQTT_PASSWORD", ""),
			QoS:            byte(getEnvAsInt("MQTT_QOS", 1)),
			CleanSession:   getEnvAsBool("MQTT_CLEAN_SESSION", true),
			ConnectTimeout: getEnvAsDuration("MQTT_CONNECT_TIMEOUT", 30*time.Second),
		},
		MonitoringManager: MonitoringManagerConfig{
			Host: getEnv("MONITORING_MANAGER_HOST", "monitoring_manager"),
			Port: getEnvAsInt("MONITORING_MANAGER_PORT", 10999),
		},
		ServiceManager: ServiceManagerConfig{
			Host: getEnv("SERVICE_MANAGER_HOST", "cluster_service_manager"),
			Port: getEnvAsInt("SERVICE_MANAGER_PORT", 10110),
		},
		MongoDB: MongoDBConfig{
			Host:     getEnv("MONGODB_HOST", "cluster_mongo_net"),
			Port:     getEnvAsInt("MONGODB_PORT", 10108),
			Username: getEnv("MONGODB_USERNAME", ""),
			Password: getEnv("MONGODB_PASSWORD", ""),
			Timeout:  getEnvAsDuration("MONGODB_TIMEOUT", 10*time.Second),
		},
	}

	// Validate configuration
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsBool gets an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsDuration gets an environment variable as a duration or returns a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
