package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smnzlnsk/routing-manager/config"
	"github.com/smnzlnsk/routing-manager/internal/api/v1/router"
	"github.com/smnzlnsk/routing-manager/internal/db/mongodb"
	"github.com/smnzlnsk/routing-manager/internal/logger"
	"github.com/smnzlnsk/routing-manager/internal/mqtt"
	"github.com/smnzlnsk/routing-manager/internal/observer/implementations"
	mongoRepo "github.com/smnzlnsk/routing-manager/internal/repository/mongodb"
	"github.com/smnzlnsk/routing-manager/internal/service"
	"github.com/smnzlnsk/routing-manager/internal/storage/memory"
	"go.uber.org/zap"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "", "Path to configuration file (YAML)")
	envFile := flag.String("env-file", "", "Path to .env file")
	useEnv := flag.Bool("use-env", false, "Use environment variables for configuration")
	logFormat := flag.String("log-format", "console", "Log format: 'console' or 'json'")
	logLevel := flag.String("log-level", "info", "Log level: 'debug', 'info', 'warn', 'error'")
	flag.Parse()

	// Initialize logger with custom configuration
	logConfig := logger.DefaultConfig()
	logConfig.Format = *logFormat
	logConfig.Level = *logLevel
	logger.Init(logConfig)
	defer logger.Sync()

	logger.Infof("Version: %s, Commit: %s, Date: %s", version, commit, date)

	// Load configuration based on flags
	var cfg *config.Config
	var err error

	// Create config factory
	configFactory := config.NewConfigLoaderFactory()

	if *configFile != "" {
		// Load from specific YAML file
		logger.Infof("Loading configuration from file: %s", *configFile)
		cfg, err = configFactory.CreateWithPath(config.YamlLoader, *configFile).Load()
	} else if *useEnv {
		// Load from environment variables
		logger.Info("Loading configuration from environment variables")
		cfg, err = configFactory.CreateWithPath(config.EnvLoader, *envFile).Load()
	} else {
		// Try automatic detection (YAML first, then environment)
		logger.Info("Attempting to load configuration from default locations")
		cfg, err = configFactory.Create(config.EnvLoader).Load()
	}

	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Log configuration source
	logger.Info("Configuration loaded successfully")

	// Setup MQTT client
	mqtt.InitInstance(config.MQTTConfig{
		Host:           cfg.MQTT.Host,
		Port:           cfg.MQTT.Port,
		ClientID:       cfg.MQTT.ClientID,
		Username:       cfg.MQTT.Username,
		Password:       cfg.MQTT.Password,
		QoS:            cfg.MQTT.QoS,
		CleanSession:   cfg.MQTT.CleanSession,
		ConnectTimeout: cfg.MQTT.ConnectTimeout,
	})

	mqttClient := mqtt.Instance()

	// Connect to MQTT broker
	if err := mqttClient.Connect(); err != nil {
		logger.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	defer mqttClient.Disconnect()

	// Initialize MongoDB connection
	mongoClient, err := mongodb.NewClient(&cfg.MongoDB, logger.Get().Desugar())
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer mongoClient.Close(ctx)

	// Setup HTTP server and services
	services, server := httpServerSetup(cfg, mongoClient)

	// Initialize observers for the interest state changes
	setupObservers(cfg, services, logger.Get().Desugar())

	go func() {
		logger.Infof("Starting server on port %d", cfg.HTTPServer.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Create storage
	store := memory.NewMemoryStore()
	defer store.Close()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-sigCh
	logger.Infof("Received signal %v, shutting down...", sig)

	// Perform graceful shutdown of services
	services.GracefulShutdown(ctx, logger.Get().Desugar())

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Failed to shutdown server: %v", err)
	}

	logger.Info("Shutdown complete")
}

// setupObservers initializes and registers observers for interest state changes
func setupObservers(cfg *config.Config, services *service.Services, logger *zap.Logger) {
	// Create task executor for the monitoring-manager
	// The service URL should ideally come from configuration
	taskExecutor := service.NewExternalTaskExecutor(
		fmt.Sprintf("http://%s:%d", cfg.MonitoringManager.Host, cfg.MonitoringManager.Port),
		5*time.Second, // Timeout
		logger,
	)

	// Create a task scheduler observer that will fire tasks to the external service
	taskSchedulerObserver := implementations.NewTaskSchedulerObserver(
		logger,
		taskExecutor,
		30*time.Second, // Task execution interval - adjust as needed
	)

	// Register observers with the subject
	services.InterestSubject.Register(taskSchedulerObserver)

	logger.Info("Interest observers registered successfully")

	// Store the task scheduler observer for graceful shutdown
	services.TaskSchedulerObserver = taskSchedulerObserver
}

func httpServerSetup(cfg *config.Config, mongoClient *mongodb.Client) (*service.Services, *http.Server) {
	// Create repositories
	repositories := mongoRepo.New(
		mongoClient.GetDatabase(),
		cfg.MongoDB.Collection,
		logger.Get().Desugar(),
	)

	// Create services
	services := service.New(repositories, logger.Get().Desugar())

	r := router.Setup(services, logger.Get().Desugar())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPServer.Port),
		Handler: r,
	}

	return services, server
}
