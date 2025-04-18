package mongodb

import (
	"github.com/smnzlnsk/routing-manager/config"
	"github.com/smnzlnsk/routing-manager/internal/db/mongodb"
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.uber.org/zap"
)

// New creates a new Repositories instance with MongoDB implementations
func New(cfg *config.MongoDBConfig, mongoClient *mongodb.Client, logger *zap.Logger) *repository.Repositories {
	return &repository.Repositories{
		AlertRepository:    NewAlertRepository(mongoClient.GetDatabase("routing"), "alerts", logger),
		InterestRepository: NewInterestRepository(mongoClient.GetDatabase("routing"), "interests", logger),
		JobRepository:      NewJobRepository(mongoClient.GetDatabase("jobs"), "jobs", logger),
		RoutingRepository:  NewRoutingRepository(mongoClient.GetDatabase("jobs"), "jobs", logger),

		// TODO: Initialize other repositories here with their dependencies
	}
}
