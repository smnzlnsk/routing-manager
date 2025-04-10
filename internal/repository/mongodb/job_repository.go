package mongodb

import (
	"context"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// jobRepository implements repository.ServiceRepository using MongoDB
type jobRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewJobRepository creates a new MongoDB-based job repository
func NewJobRepository(db *mongo.Database, collection string, logger *zap.Logger) repository.JobRepository {
	coll := db.Collection(collection)

	return &jobRepository{
		collection: coll,
		logger:     logger,
	}
}

func (r *jobRepository) GetByJobName(ctx context.Context, jobName string) (*domain.Job, error) {
	r.logger.Debug("Getting job by job name from MongoDB", zap.String("jobName", jobName))

	var job domain.Job
	err := r.collection.FindOne(ctx, bson.M{"job_name": jobName}).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &job, nil
}
