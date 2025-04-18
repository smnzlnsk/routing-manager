package mongodb

import (
	"context"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
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

func (r *routingRepository) GetRouting(ctx context.Context, jobName string) (*domain.JobRouting, error) {
	r.logger.Debug("Getting routing for job", zap.String("jobName", jobName))

	filter := bson.M{"job_name": jobName}
	var routing domain.JobRouting
	err := r.collection.FindOne(ctx, filter).Decode(&routing)
	if err != nil {
		return nil, err
	}
	return &routing, nil
}

func (r *routingRepository) UpdateRouting(ctx context.Context, routing *domain.JobRouting) error {
	r.logger.Debug("Updating routing priorities for job", zap.String("jobName", routing.JobName))

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"job_name": routing.JobName},
		bson.M{"$set": bson.M{"service_ip_priority": routing.ServiceIPPriority}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		r.logger.Debug("No document found with the specified job name", zap.String("jobName", routing.JobName))
		return domain.ErrNotFound
	}
	return nil
}
