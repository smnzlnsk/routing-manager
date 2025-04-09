package mongodb

import (
	"context"
	"time"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// interestRepository implements repository.InterestRepository using MongoDB
type interestRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewInterestRepository creates a new MongoDB-based interest repository
func NewInterestRepository(db *mongo.Database, collection string, logger *zap.Logger) repository.InterestRepository {
	coll := db.Collection(collection)

	// Create unique index on AppName
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "appname", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	// Create a non-unique index for ServiceIp for efficient querying
	serviceIpIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "serviceip", Value: 1},
		},
		// Non-unique index (no SetUnique)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := coll.Indexes().CreateOne(ctx, indexModel); err != nil {
		logger.Error("Failed to create index on appname", zap.Error(err))
	}

	if _, err := coll.Indexes().CreateOne(ctx, serviceIpIndex); err != nil {
		logger.Error("Failed to create index on serviceip", zap.Error(err))
	}

	return &interestRepository{
		collection: coll,
		logger:     logger,
	}
}

// Create adds a new interest to the database
func (r *interestRepository) Create(ctx context.Context, interest *domain.Interest) error {
	r.logger.Debug("Creating interest in MongoDB", zap.String("appName", interest.AppName))

	// Convert domain.Interest to BSON document
	doc := bson.M{
		"appname":   interest.AppName,
		"serviceip": interest.ServiceIp,
		"createdat": interest.CreatedAt,
		"updatedat": interest.UpdatedAt,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		// Check if the error is a duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrInterestAlreadyExists
		}
		return err
	}

	return nil
}

// GetByAppName retrieves an interest by its app name
func (r *interestRepository) GetByAppName(ctx context.Context, appName string) (*domain.Interest, error) {
	r.logger.Debug("Getting interest by app name from MongoDB", zap.String("appName", appName))

	var result bson.M
	err := r.collection.FindOne(ctx, bson.M{"appname": appName}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	// Convert BSON to domain.Interest
	interest := &domain.Interest{
		AppName:   result["appname"].(string),
		ServiceIp: result["serviceip"].(string),
	}

	// Handle timestamps
	if createdAt, ok := result["createdat"].(primitive.DateTime); ok {
		interest.CreatedAt = createdAt.Time()
	}
	if updatedAt, ok := result["updatedat"].(primitive.DateTime); ok {
		interest.UpdatedAt = updatedAt.Time()
	}

	return interest, nil
}

// GetByServiceIp retrieves an interest by its service IP
func (r *interestRepository) GetByServiceIp(ctx context.Context, serviceIp string) (*domain.Interest, error) {
	r.logger.Debug("Getting interest by service IP from MongoDB", zap.String("serviceIp", serviceIp))

	var result bson.M
	err := r.collection.FindOne(ctx, bson.M{"serviceip": serviceIp}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	// Convert BSON to domain.Interest
	interest := &domain.Interest{
		AppName:   result["appname"].(string),
		ServiceIp: result["serviceip"].(string),
	}

	// Handle timestamps
	if createdAt, ok := result["createdat"].(primitive.DateTime); ok {
		interest.CreatedAt = createdAt.Time()
	}
	if updatedAt, ok := result["updatedat"].(primitive.DateTime); ok {
		interest.UpdatedAt = updatedAt.Time()
	}

	return interest, nil
}

// Update updates an existing interest
func (r *interestRepository) Update(ctx context.Context, interest *domain.Interest) (*domain.Interest, error) {
	r.logger.Debug("Updating interest in MongoDB", zap.String("appName", interest.AppName))

	// Prepare update document
	update := bson.M{
		"$set": bson.M{
			"serviceip": interest.ServiceIp,
			"updatedat": time.Now(),
		},
	}

	result := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"appname": interest.AppName},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, result.Err()
	}

	var updatedDoc bson.M
	if err := result.Decode(&updatedDoc); err != nil {
		return nil, err
	}

	// Convert BSON to domain.Interest
	updatedInterest := &domain.Interest{
		AppName:   updatedDoc["appname"].(string),
		ServiceIp: updatedDoc["serviceip"].(string),
	}

	// Handle timestamps
	if createdAt, ok := updatedDoc["createdat"].(primitive.DateTime); ok {
		updatedInterest.CreatedAt = createdAt.Time()
	}
	if updatedAt, ok := updatedDoc["updatedat"].(primitive.DateTime); ok {
		updatedInterest.UpdatedAt = updatedAt.Time()
	}

	return updatedInterest, nil
}

// DeleteByAppName deletes an interest by its app name
func (r *interestRepository) DeleteByAppName(ctx context.Context, appName string) error {
	r.logger.Debug("Deleting interest by app name from MongoDB", zap.String("appName", appName))

	result, err := r.collection.DeleteOne(ctx, bson.M{"appname": appName})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// DeleteByServiceIp deletes an interest by its service IP
func (r *interestRepository) DeleteByServiceIp(ctx context.Context, serviceIp string) error {
	r.logger.Debug("Deleting interest by service IP from MongoDB", zap.String("serviceIp", serviceIp))

	result, err := r.collection.DeleteOne(ctx, bson.M{"serviceip": serviceIp})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// List retrieves all interests
func (r *interestRepository) List(ctx context.Context) ([]*domain.Interest, error) {
	r.logger.Debug("Listing all interests from MongoDB")

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var interests []*domain.Interest
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		interest := &domain.Interest{
			AppName:   result["appname"].(string),
			ServiceIp: result["serviceip"].(string),
		}

		// Handle timestamps
		if createdAt, ok := result["createdat"].(primitive.DateTime); ok {
			interest.CreatedAt = createdAt.Time()
		}
		if updatedAt, ok := result["updatedat"].(primitive.DateTime); ok {
			interest.UpdatedAt = updatedAt.Time()
		}

		interests = append(interests, interest)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return interests, nil
}
