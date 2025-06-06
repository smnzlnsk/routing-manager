package service

import (
	"github.com/smnzlnsk/routing-manager/internal/observer"
	"github.com/smnzlnsk/routing-manager/internal/observer/implementations"
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.uber.org/zap"
)

// Services is a collection of all services in the application
type Services struct {
	AlertService          AlertService
	InterestService       InterestService
	InterestSubject       *observer.InterestSubject
	TaskSchedulerObserver *implementations.TaskSchedulerObserver
	JobService            JobService
	RoutingService        RoutingService
}

// NewServices creates a new Services instance
func New(repositories *repository.Repositories, logger *zap.Logger) *Services {
	// Create the interest subject for observer pattern
	interestSubject := observer.NewInterestSubject(logger)

	return &Services{
		AlertService:    NewAlertService(repositories.AlertRepository, logger),
		InterestService: NewInterestService(repositories.InterestRepository, interestSubject, logger),
		InterestSubject: interestSubject,
		// TaskSchedulerObserver will be set separately after creation
		JobService:     NewJobService(repositories.JobRepository, logger),
		RoutingService: NewRoutingService(repositories.RoutingRepository, logger),
		// Initialize other services here with their dependencies
	}
}
