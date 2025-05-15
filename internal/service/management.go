package service

import (
	"context"
	"sync"
	"time"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"go.uber.org/zap"
)

var restartOnce sync.Once

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

// Restart is supposed to be run once, in case of a restart of the routing-manager due to a crash
func (s *Services) Restart(ctx context.Context, logger *zap.Logger) {
	restartOnce.Do(func() {
		logger.Info("Starting service restart procedure")

		// Check if we have necessary components
		if s.InterestService == nil || s.InterestSubject == nil || s.TaskSchedulerObserver == nil {
			logger.Error("Cannot restart services: required components are missing")
			return
		}

		// Fetch all interests from the database
		interests, err := s.InterestService.List(ctx)
		if err != nil {
			logger.Error("Failed to retrieve interests during restart", zap.Error(err))
			return
		}

		logger.Info("Retrieved interests from database", zap.Int("count", len(interests)))

		// First make sure any existing tasks are stopped
		s.TaskSchedulerObserver.Shutdown()

		// Reinitialize by notifying about all existing interests
		for _, interest := range interests {
			// Create an interest event and notify the subject
			event := domain.InterestEvent{
				Type:     domain.InterestCreated,
				Interest: interest,
			}

			// Notify the subject to trigger observer updates
			s.InterestSubject.Notify(event)

			logger.Debug("Reinitialized interest",
				zap.String("appName", interest.AppName),
				zap.String("serviceIp", interest.ServiceIp))
		}

		logger.Info("Service restart procedure completed successfully", zap.Int("interestsReinitialized", len(interests)))
	})
}
