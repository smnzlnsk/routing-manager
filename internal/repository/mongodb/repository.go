package mongodb

import (
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// New creates a new Repositories instance with MongoDB implementations
func New(db *mongo.Database, collectionName string, logger *zap.Logger) *repository.Repositories {
	return &repository.Repositories{
		InterestRepository: NewInterestRepository(db, collectionName, logger),
		// Initialize other repositories here with their dependencies
	}
}
