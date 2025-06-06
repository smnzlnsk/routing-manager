package mongodb

import (
	"context"
	"fmt"

	"github.com/smnzlnsk/routing-manager/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Client represents a MongoDB client with connection to a specific database
type Client struct {
	client    *mongo.Client
	databases map[string]*mongo.Database
	logger    *zap.Logger
}

// NewClient creates a new MongoDB client
func NewClient(cfg *config.MongoDBConfig, logger *zap.Logger) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// Set up MongoDB client options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port))

	// Set authentication credentials if provided
	if cfg.Username != "" && cfg.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: cfg.Username,
			Password: cfg.Password,
		})
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	// Connect to the databases
	databases := make(map[string]*mongo.Database)
	databases["routing"] = client.Database("routing")
	databases["jobs"] = client.Database("jobs")

	logger.Info("Connected to MongoDB",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Any("databases", databases))

	return &Client{
		client:    client,
		databases: databases,
		logger:    logger,
	}, nil
}

// GetDatabase returns the MongoDB database
func (c *Client) GetDatabase(database string) *mongo.Database {
	return c.databases[database]
}

// Close closes the MongoDB connection
func (c *Client) Close(ctx context.Context) error {
	c.logger.Info("Closing MongoDB connection")
	return c.client.Disconnect(ctx)
}
