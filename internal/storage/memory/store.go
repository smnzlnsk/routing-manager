package memory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/smnzlnsk/routing-manager/internal/storage"
)

var (
	// ErrServiceNotFound is returned when a service is not found in the store
	ErrServiceNotFound = errors.New("service not found")
)

// MemoryStore implements the storage.PerformanceStore interface using in-memory storage
type MemoryStore struct {
	mu sync.RWMutex
	// Map of service ID to its latest performance metric
	latestMetrics map[string]float64
	// Map of service ID to its historical performance metrics
	historicalMetrics map[string][]storage.PerformanceRecord
}

// NewMemoryStore creates a new in-memory performance store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		latestMetrics:     make(map[string]float64),
		historicalMetrics: make(map[string][]storage.PerformanceRecord),
	}
}

// SaveMetric stores a performance metric for a service
func (m *MemoryStore) SaveMetric(ctx context.Context, serviceID string, value float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Store the latest metric
	m.latestMetrics[serviceID] = value

	// Store in historical metrics
	record := storage.PerformanceRecord{
		ServiceID:  serviceID,
		Value:      value,
		RecordedAt: time.Now(),
	}

	if _, exists := m.historicalMetrics[serviceID]; !exists {
		m.historicalMetrics[serviceID] = []storage.PerformanceRecord{}
	}

	m.historicalMetrics[serviceID] = append(m.historicalMetrics[serviceID], record)

	return nil
}

// UpdateMetric updates an existing performance metric
func (m *MemoryStore) UpdateMetric(ctx context.Context, serviceID string, value float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the service exists
	if _, exists := m.latestMetrics[serviceID]; !exists {
		return ErrServiceNotFound
	}

	// Update the latest metric
	m.latestMetrics[serviceID] = value

	// Add to historical metrics
	record := storage.PerformanceRecord{
		ServiceID:  serviceID,
		Value:      value,
		RecordedAt: time.Now(),
	}

	m.historicalMetrics[serviceID] = append(m.historicalMetrics[serviceID], record)

	return nil
}

// GetMetric retrieves the latest performance metric for a service
func (m *MemoryStore) GetMetric(ctx context.Context, serviceID string) (float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.latestMetrics[serviceID]
	if !exists {
		return 0, ErrServiceNotFound
	}

	return value, nil
}

// GetMetricHistory retrieves historical performance metrics for a service
func (m *MemoryStore) GetMetricHistory(ctx context.Context, serviceID string, since time.Time, limit int) ([]storage.PerformanceRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	records, exists := m.historicalMetrics[serviceID]
	if !exists {
		return nil, ErrServiceNotFound
	}

	// Filter records by time
	var filteredRecords []storage.PerformanceRecord
	for _, record := range records {
		if record.RecordedAt.After(since) || record.RecordedAt.Equal(since) {
			filteredRecords = append(filteredRecords, record)
		}
	}

	// Sort by time (newest first)
	sort.Slice(filteredRecords, func(i, j int) bool {
		return filteredRecords[i].RecordedAt.After(filteredRecords[j].RecordedAt)
	})

	// Apply limit if specified
	if limit > 0 && len(filteredRecords) > limit {
		filteredRecords = filteredRecords[:limit]
	}

	return filteredRecords, nil
}

// GetAllMetrics retrieves the latest performance metrics for all services
func (m *MemoryStore) GetAllMetrics(ctx context.Context) (map[string]float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy of the metrics map to avoid concurrent access issues
	result := make(map[string]float64, len(m.latestMetrics))
	for serviceID, value := range m.latestMetrics {
		result[serviceID] = value
	}

	return result, nil
}

// Close closes the storage connection (no-op for in-memory store)
func (m *MemoryStore) Close() error {
	return nil
}
