package service

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// GracefulShutdown performs a graceful shutdown of all services
func (s *Services) GracefulShutdown(ctx context.Context, logger *zap.Logger) {
	// Set a timeout for the shutdown
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	logger.Info("Starting graceful shutdown of services")

	// Shutdown task scheduler observer if it exists
	if s.TaskSchedulerObserver != nil {
		s.TaskSchedulerObserver.Shutdown()
		logger.Info("Task scheduler observer shut down successfully")
	}

	// Add other service shutdown logic here
	// ...

	logger.Info("All services shut down successfully")
}
