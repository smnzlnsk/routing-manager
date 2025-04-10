package mongodb

import (
	"context"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// alertRepository implements repository.AlertRepository using MongoDB
type alertRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewAlertRepository creates a new MongoDB-based alert repository
func NewAlertRepository(db *mongo.Database, collection string, logger *zap.Logger) repository.AlertRepository {
	coll := db.Collection(collection)

	return &alertRepository{
		collection: coll,
		logger:     logger,
	}
}

// Create adds a new alert to the database
func (r *alertRepository) Create(ctx context.Context, alert *domain.Alert) error {
	r.logger.Debug("Creating alert in MongoDB", zap.String("appName", alert.AppName))

	_, err := r.collection.InsertOne(ctx, alert)
	if err != nil {
		return err
	}
	return nil
}

// GetByAppName retrieves an alert by app name
func (r *alertRepository) GetByAppName(ctx context.Context, appName string) (*domain.Alert, error) {
	r.logger.Debug("Getting alert by app name", zap.String("appName", appName))

	// Find the alert by app name
	var alert domain.Alert
	err := r.collection.FindOne(ctx, bson.M{"appname": appName}).Decode(&alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}
