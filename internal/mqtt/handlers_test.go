package mqtt

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/smnzlnsk/routing-manager/internal/logger"
	"github.com/stretchr/testify/assert"
)

// TestHandleDeviceDataMessage tests the device data message handler
func TestHandleDeviceDataMessage(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name      string
		topic     string
		variables map[string]string
		payload   []byte
		wantErr   bool
	}{
		{
			name:  "Valid device data",
			topic: "devices/device123/sensors/temp1/data",
			variables: map[string]string{
				"deviceId": "device123",
				"sensorId": "temp1",
			},
			payload: []byte(`{"temperature": 22.5, "humidity": 45}`),
			wantErr: false,
		},
		{
			name:  "Invalid JSON payload",
			topic: "devices/device123/sensors/temp1/data",
			variables: map[string]string{
				"deviceId": "device123",
				"sensorId": "temp1",
			},
			payload: []byte(`{"temperature": 22.5, "humidity": 45`), // Invalid JSON
			wantErr: true,
		},
		{
			name:  "Missing device ID",
			topic: "devices/sensors/temp1/data",
			variables: map[string]string{
				"sensorId": "temp1",
			},
			payload: []byte(`{"temperature": 22.5, "humidity": 45}`),
			wantErr: false, // Should not error, but deviceId will be empty
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the handler
			handleDeviceDataMessage(tt.topic, tt.variables, tt.payload)

			// Since we can't easily mock the logger in this test,
			// we're just verifying that the function doesn't panic
			// A more comprehensive test would mock the logger and verify outputs
		})
	}
}

// TestHandleWildcardTopic tests the generic wildcard topic handler
func TestHandleWildcardTopic(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name         string
		topicPattern string
		topic        string
		variables    map[string]string
		payload      []byte
	}{
		{
			name:         "User status topic",
			topicPattern: "users/{userId}/status",
			topic:        "users/user123/status",
			variables: map[string]string{
				"userId": "user123",
			},
			payload: []byte(`{"status": "online", "lastSeen": "2023-04-01T12:00:00Z"}`),
		},
		{
			name:         "Multiple variables",
			topicPattern: "{region}/devices/{deviceId}/commands",
			topic:        "eu-west/devices/device456/commands",
			variables: map[string]string{
				"region":   "eu-west",
				"deviceId": "device456",
			},
			payload: []byte(`{"command": "restart", "delay": 5}`),
		},
		{
			name:         "No variables",
			topicPattern: "broadcast/all",
			topic:        "broadcast/all",
			variables:    map[string]string{},
			payload:      []byte(`{"message": "System maintenance in 10 minutes"}`),
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the handler
			handler := handleWildcardTopic(tt.topicPattern)

			// Call the handler
			handler(tt.topic, tt.variables, tt.payload)

			// Since we can't easily mock the logger in this test,
			// we're just verifying that the function doesn't panic
			// A more comprehensive test would mock the logger and verify outputs
		})
	}
}

// TestIntegrationWithTopicParser tests the integration between the topic parser and handlers
func TestIntegrationWithTopicParser(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name         string
		topicPattern string
		actualTopic  string
		payload      []byte
		expectedVars map[string]string
	}{
		{
			name:         "Device data topic",
			topicPattern: "devices/{deviceId}/sensors/{sensorId}/data",
			actualTopic:  "devices/thermostat1/sensors/temp/data",
			payload:      []byte(`{"temperature": 22.5}`),
			expectedVars: map[string]string{
				"deviceId": "thermostat1",
				"sensorId": "temp",
			},
		},
		{
			name:         "User activity topic",
			topicPattern: "users/{userId}/activities/{activityType}",
			actualTopic:  "users/user123/activities/login",
			payload:      []byte(`{"timestamp": "2023-04-01T12:00:00Z", "ip": "192.168.1.1"}`),
			expectedVars: map[string]string{
				"userId":       "user123",
				"activityType": "login",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a topic parser
			parser := newTopicParser(tt.topicPattern)

			// Extract variables
			variables := parser.extractVariables(tt.actualTopic)

			// Verify extracted variables match expected
			assert.Equal(t, tt.expectedVars, variables, "Extracted variables should match expected")

			// Test with device data handler if it's a device data topic
			if tt.topicPattern == "devices/{deviceId}/sensors/{sensorId}/data" {
				// This just verifies it doesn't panic
				handleDeviceDataMessage(tt.actualTopic, variables, tt.payload)
			} else {
				// Use the generic handler for other patterns
				handler := handleWildcardTopic(tt.topicPattern)
				handler(tt.actualTopic, variables, tt.payload)
			}
		})
	}
}

// TestProcessDeviceData tests the device data processing function
func TestProcessDeviceData(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name     string
		deviceID string
		sensorID string
		data     map[string]interface{}
	}{
		{
			name:     "Temperature data",
			deviceID: "thermostat1",
			sensorID: "temp",
			data: map[string]interface{}{
				"temperature": 22.5,
				"unit":        "celsius",
			},
		},
		{
			name:     "Humidity data",
			deviceID: "humidifier2",
			sensorID: "humidity",
			data: map[string]interface{}{
				"humidity": 45.0,
				"unit":     "percent",
			},
		},
		{
			name:     "Empty data",
			deviceID: "device3",
			sensorID: "status",
			data:     map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function
			processDeviceData(tt.deviceID, tt.sensorID, tt.data)

			// Since we can't easily mock the logger in this test,
			// we're just verifying that the function doesn't panic
			// A more comprehensive test would mock the logger and verify outputs
		})
	}
}

// TestMQTTWildcardPatterns tests various MQTT wildcard patterns and their variable extraction
func TestMQTTWildcardPatterns(t *testing.T) {
	tests := []struct {
		name         string
		topicPattern string
		topics       map[string]map[string]string // Map of topic to expected variables
	}{
		{
			name:         "Single-level wildcard (+)",
			topicPattern: "devices/{deviceId}/+/{property}",
			topics: map[string]map[string]string{
				"devices/device123/temperature/value": {
					"deviceId": "device123",
					"property": "value",
				},
				"devices/device456/humidity/status": {
					"deviceId": "device456",
					"property": "status",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := newTopicParser(tt.topicPattern)

			for topic, expectedVars := range tt.topics {
				variables := parser.extractVariables(topic)

				// Check if the extracted variables match the expected ones
				for key, expectedValue := range expectedVars {
					assert.Contains(t, variables, key, "Variable should be extracted")
					assert.Equal(t, expectedValue, variables[key], "Variable value should match")
				}

				// Create a handler for this pattern
				handler := handleWildcardTopic(tt.topicPattern)

				// Verify the handler doesn't panic
				handler(topic, variables, []byte(`{"test": true}`))
			}
		})
	}
}

// HandleDeviceDataMessage handles messages from device data topics with variables
// Example topic pattern: devices/{deviceId}/sensors/{sensorId}/data
func handleDeviceDataMessage(_ string, variables map[string]string, payload []byte) {
	deviceID := variables["deviceId"]
	sensorID := variables["sensorId"]

	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err == nil {
		processDeviceData(deviceID, sensorID, data)
	} else {
		fmt.Printf("Failed to parse device data payload: %v", err)
	}
}

// processDeviceData is an example function that processes device data
func processDeviceData(_, _ string, _ map[string]interface{}) {}

// HandleWildcardTopic is a generic handler for any topic with wildcards
func handleWildcardTopic(topicPattern string) WildcardMessageHandler {
	return func(topic string, variables map[string]string, payload []byte) {
		// Log the received message with all variables
		logger.Infof("Received message on wildcard topic",
			"topic", topic, "\n",
			"pattern", topicPattern, "\n",
			"payload", string(payload))

		// Add all variables to log fields as separate log entries
		for key, value := range variables {
			logger.Infof("Topic variable: %s = %s", key, value)
		}

		// Process the message based on the variables
		// This is where you would implement your custom logic
	}
}
