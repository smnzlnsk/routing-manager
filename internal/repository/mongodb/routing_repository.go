package mongodb

import (
	"context"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// routingRepository implements repository.RoutingRepository using MongoDB
type routingRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewRoutingRepository creates a new MongoDB-based routing repository
// It holds the collection of current routing priorities for each service
func NewRoutingRepository(db *mongo.Database, collection string, logger *zap.Logger) repository.RoutingRepository {
	coll := db.Collection(collection)

	return &routingRepository{
		collection: coll,
		logger:     logger,
	}
}

func (r *routingRepository) GetRouting(ctx context.Context, jobName string) (*domain.Job, error) {
	r.logger.Debug("Getting routing for job", zap.String("jobName", jobName))
	return &domain.Job{}, nil
}

func (r *routingRepository) UpdateRouting(ctx context.Context, routing *domain.Job) error {
	r.logger.Debug("Updating routing priorities for job", zap.String("jobName", routing.JobName))
	return nil
}
