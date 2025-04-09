package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	observer "github.com/smnzlnsk/routing-manager/internal/observer/implementations"
	"go.uber.org/zap"
)

// ExternalTaskExecutor implements the TaskExecutor interface for sending tasks to external microservices
type ExternalTaskExecutor struct {
	httpClient *http.Client
	serviceURL string
	logger     *zap.Logger
}

// TaskPayload represents the data to be sent to the external service
type TaskPayload struct {
	AppName     string    `json:"appName"`
	ServiceIP   string    `json:"serviceIp"`
	Timestamp   time.Time `json:"timestamp"`
	RequestType string    `json:"requestType"`
}

// NewExternalTaskExecutor creates a new instance of ExternalTaskExecutor
func NewExternalTaskExecutor(serviceURL string, timeout time.Duration, logger *zap.Logger) observer.TaskExecutor {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &ExternalTaskExecutor{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		serviceURL: serviceURL,
		logger:     logger,
	}
}

// ExecuteTask sends a task request to the external microservice
func (e *ExternalTaskExecutor) ExecuteTask(interest *domain.Interest) error {
	payload := TaskPayload{
		AppName:     interest.AppName,
		ServiceIP:   interest.ServiceIp,
		Timestamp:   time.Now(),
		RequestType: "health_check", // Can be parameterized if needed
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	// Construct the target URL
	targetURL := fmt.Sprintf("%s/policy/routing/def", e.serviceURL)

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set appropriate headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Source", "routing-manager")

	// Execute the request
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute task request: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("task request failed with status code: %d", resp.StatusCode)
	}

	e.logger.Info("Task executed successfully",
		zap.String("appName", interest.AppName),
		zap.String("serviceIP", interest.ServiceIp),
		zap.Int("statusCode", resp.StatusCode),
		zap.String("body", string(respBody)))

	return nil
}
