package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/service"
	"go.uber.org/zap"
)

// ExternalTaskExecutor implements the TaskExecutor interface for sending tasks to external microservices
type ExternalTaskExecutor struct {
	httpClient *http.Client
	serviceURL string
	logger     *zap.Logger
	jobService service.JobService
}

// TaskPayload represents the data to be sent to the external service
type TaskPayload struct {
	AppName   string                 `json:"appName"`
	ServiceIP string                 `json:"serviceIp"`
	Timestamp time.Time              `json:"timestamp"`
	JobData   map[string]interface{} `json:"jobData,omitempty"`
}

// NewExternalTaskExecutor creates a new instance of ExternalTaskExecutor
func NewExternalTaskExecutor(serviceURL string, timeout time.Duration, jobService service.JobService, logger *zap.Logger) domain.TaskExecutor {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &ExternalTaskExecutor{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		serviceURL: serviceURL,
		jobService: jobService,
		logger:     logger,
	}
}

// ExecuteTask sends a task request to the external microservice
func (e *ExternalTaskExecutor) ExecuteTask(interest *domain.Interest) error {
	// Create a basic payload
	payload := TaskPayload{
		AppName:   interest.AppName,
		ServiceIP: interest.ServiceIp,
		Timestamp: time.Now(),
	}

	// If we need job data, retrieve it
	if job, err := e.jobService.GetByJobName(context.Background(), interest.AppName); err == nil {
		// We found a job, add its data to the payload
		jobData := make(map[string]interface{})
		// Add job data as needed
		jobData["job_name"] = job.JobName
		jobData["service_ip_list"] = job.ServiceIpList
		jobData["instance_list"] = job.ServiceInstanceList
		// Add the job data to the payload
		payload.JobData = jobData
	} else {
		e.logger.Warn("Could not find job data for interest",
			zap.String("appName", interest.AppName),
			zap.Error(err))
	}

	for _, entry := range payload.JobData["service_ip_list"].([]domain.ServiceIpListEntry) {
		ipType := entry.IpType

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal task payload: %w", err)
		}

		// Construct the target URL
		targetURL := fmt.Sprintf("%s/policy/routing/%s", e.serviceURL, ipType)

		// Create the HTTP request
		req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Set appropriate headers
		req.Header.Set("Content-Type", "application/json")

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
			return fmt.Errorf("task request failed with status code: %d, body: %s", resp.StatusCode, string(respBody))
		}

		e.logger.Info("Task executed successfully",
			zap.String("appName", interest.AppName),
			zap.String("serviceIP", interest.ServiceIp),
			zap.Int("statusCode", resp.StatusCode),
			zap.String("body", string(respBody)))
	}

	return nil
}
