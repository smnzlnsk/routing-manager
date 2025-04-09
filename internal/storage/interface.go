package storage

import (
	"context"
	"time"
)

// PerformanceStore defines the interface for storing and retrieving
// performance metrics for services
type PerformanceStore interface {
	// SaveMetric stores a performance metric for a service
	SaveMetric(ctx context.Context, serviceID string, value float64) error

	// UpdateMetric updates an existing performance metric
	UpdateMetric(ctx context.Context, serviceID string, value float64) error

	// GetMetric retrieves the latest performance metric for a service
	GetMetric(ctx context.Context, serviceID string) (float64, error)

	// GetMetricHistory retrieves historical performance metrics for a service
	GetMetricHistory(ctx context.Context, serviceID string, since time.Time, limit int) ([]PerformanceRecord, error)

	// GetAllMetrics retrieves the latest performance metrics for all services
	GetAllMetrics(ctx context.Context) (map[string]float64, error)

	// Close closes the storage connection
	Close() error
}

// PerformanceRecord represents a single performance metric record with timestamp
type PerformanceRecord struct {
	ServiceID  string    `json:"service_id"`
	Value      float64   `json:"value"`
	RecordedAt time.Time `json:"recorded_at"`
}
