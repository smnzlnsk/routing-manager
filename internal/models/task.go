package models

import (
	"time"
)

// Task represents a task received from MQTT that needs to be processed
type Task struct {
	ID          string                 `json:"id"`
	ServiceID   string                 `json:"service_id"`
	Action      string                 `json:"action"`
	Payload     map[string]interface{} `json:"payload"`
	ReceivedAt  time.Time              `json:"received_at"`
	ProcessedAt time.Time              `json:"processed_at,omitempty"`
}

// TaskResult represents the result of a task after processing
type TaskResult struct {
	TaskID       string    `json:"task_id"`
	ServiceID    string    `json:"service_id"`
	Success      bool      `json:"success"`
	Performance  float64   `json:"performance"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CompletedAt  time.Time `json:"completed_at"`
}

// PerformanceMetric represents a performance metric for a service
type PerformanceMetric struct {
	ServiceID   string    `json:"service_id"`
	Value       float64   `json:"value"`
	RecordedAt  time.Time `json:"recorded_at"`
	Description string    `json:"description,omitempty"`
	Unit        string    `json:"unit,omitempty"`
}
