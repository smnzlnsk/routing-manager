package mqtt

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTopicParser(t *testing.T) {
	tests := []struct {
		name            string
		pattern         string
		wantMQTTPattern string
		wantVarNames    []string
	}{
		{
			name:            "Simple pattern with no variables",
			pattern:         "devices/data",
			wantMQTTPattern: "devices/data",
			wantVarNames:    []string{},
		},
		{
			name:            "Pattern with one variable",
			pattern:         "devices/{deviceId}/data",
			wantMQTTPattern: "devices/+/data",
			wantVarNames:    []string{"deviceId"},
		},
		{
			name:            "Pattern with multiple variables",
			pattern:         "devices/{deviceId}/sensors/{sensorId}/data",
			wantMQTTPattern: "devices/+/sensors/+/data",
			wantVarNames:    []string{"deviceId", "sensorId"},
		},
		{
			name:            "Pattern with adjacent variables",
			pattern:         "users/{userId}/{action}",
			wantMQTTPattern: "users/+/+",
			wantVarNames:    []string{"userId", "action"},
		},
		{
			name:            "Pattern with variable at start",
			pattern:         "{region}/devices/{deviceId}",
			wantMQTTPattern: "+/devices/+",
			wantVarNames:    []string{"region", "deviceId"},
		},
		{
			name:            "Pattern with variable at end",
			pattern:         "devices/{deviceId}/{property}",
			wantMQTTPattern: "devices/+/+",
			wantVarNames:    []string{"deviceId", "property"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := newTopicParser(tt.pattern)

			if parser.mqttPattern != tt.wantMQTTPattern {
				t.Errorf("newTopicParser() mqttPattern = %v, want %v", parser.mqttPattern, tt.wantMQTTPattern)
			}

			if !reflect.DeepEqual(parser.variableNames, tt.wantVarNames) {
				t.Errorf("newTopicParser() variableNames = %v, want %v", parser.variableNames, tt.wantVarNames)
			}
		})
	}
}

func TestGetMQTTSubscriptionTopic(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    string
	}{
		{
			name:    "Simple pattern",
			pattern: "devices/data",
			want:    "devices/data",
		},
		{
			name:    "Pattern with variables",
			pattern: "devices/{deviceId}/sensors/{sensorId}/data",
			want:    "devices/+/sensors/+/data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := newTopicParser(tt.pattern)
			if got := parser.getMQTTSubscriptionTopic(); got != tt.want {
				t.Errorf("getMQTTSubscriptionTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractVariables(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		topic   string
		want    map[string]string
	}{
		{
			name:    "Simple pattern with no variables",
			pattern: "devices/data",
			topic:   "devices/data",
			want:    map[string]string{},
		},
		{
			name:    "Pattern with one variable",
			pattern: "devices/{deviceId}/data",
			topic:   "devices/device123/data",
			want:    map[string]string{"deviceId": "device123"},
		},
		{
			name:    "Pattern with multiple variables",
			pattern: "devices/{deviceId}/sensors/{sensorId}/data",
			topic:   "devices/device123/sensors/temp1/data",
			want:    map[string]string{"deviceId": "device123", "sensorId": "temp1"},
		},
		{
			name:    "Pattern with adjacent variables",
			pattern: "users/{userId}/{action}",
			topic:   "users/user456/login",
			want:    map[string]string{"userId": "user456", "action": "login"},
		},
		{
			name:    "Pattern with variable at start",
			pattern: "{region}/devices/{deviceId}",
			topic:   "eu-west/devices/device789",
			want:    map[string]string{"region": "eu-west", "deviceId": "device789"},
		},
		{
			name:    "Pattern with variable at end",
			pattern: "devices/{deviceId}/{property}",
			topic:   "devices/device123/temperature",
			want:    map[string]string{"deviceId": "device123", "property": "temperature"},
		},
		{
			name:    "Topic with special characters in variables",
			pattern: "devices/{deviceId}/data",
			topic:   "devices/device-123_456.789/data",
			want:    map[string]string{"deviceId": "device-123_456.789"},
		},
		{
			name:    "Topic with length mismatch",
			pattern: "devices/{deviceId}/sensors/{sensorId}/data",
			topic:   "devices/device123/data",
			want:    map[string]string{}, // Should return empty map due to length mismatch
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := newTopicParser(tt.pattern)
			if got := parser.extractVariables(tt.topic); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComplexScenarios(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		topics  map[string]map[string]string // map of topic to expected variables
	}{
		{
			name:    "IoT device data pattern",
			pattern: "devices/{deviceId}/sensors/{sensorId}/data",
			topics: map[string]map[string]string{
				"devices/thermostat1/sensors/temp/data": {
					"deviceId": "thermostat1",
					"sensorId": "temp",
				},
				"devices/gateway2/sensors/humidity/data": {
					"deviceId": "gateway2",
					"sensorId": "humidity",
				},
				"devices/invalid/wrong/format": {}, // Should return empty map due to format mismatch
			},
		},
		{
			name:    "User activity pattern",
			pattern: "users/{userId}/activities/{activityType}",
			topics: map[string]map[string]string{
				"users/user123/activities/login": {
					"userId":       "user123",
					"activityType": "login",
				},
				"users/admin456/activities/configuration": {
					"userId":       "admin456",
					"activityType": "configuration",
				},
			},
		},
		{
			name:    "Multi-level pattern",
			pattern: "{region}/{country}/stores/{storeId}/departments/{deptId}",
			topics: map[string]map[string]string{
				"europe/germany/stores/berlin1/departments/electronics": {
					"region":  "europe",
					"country": "germany",
					"storeId": "berlin1",
					"deptId":  "electronics",
				},
				"namerica/usa/stores/nyc3/departments/clothing": {
					"region":  "namerica",
					"country": "usa",
					"storeId": "nyc3",
					"deptId":  "clothing",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := newTopicParser(tt.pattern)

			// Verify MQTT pattern conversion
			expectedMQTTPattern := tt.pattern
			for _, varName := range parser.variableNames {
				expectedMQTTPattern = strings.Replace(expectedMQTTPattern, "{"+varName+"}", "+", 1)
			}
			if parser.getMQTTSubscriptionTopic() != expectedMQTTPattern {
				t.Errorf("getMQTTSubscriptionTopic() = %v, want %v",
					parser.getMQTTSubscriptionTopic(), expectedMQTTPattern)
			}

			// Test variable extraction for each topic
			for topic, expectedVars := range tt.topics {
				extractedVars := parser.extractVariables(topic)
				if !reflect.DeepEqual(extractedVars, expectedVars) {
					t.Errorf("For topic %s: extractVariables() = %v, want %v",
						topic, extractedVars, expectedVars)
				}
			}
		})
	}
}

// TestEdgeCases tests some edge cases for the topic parser
func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		topic   string
		want    map[string]string
	}{
		{
			name:    "Empty pattern",
			pattern: "",
			topic:   "",
			want:    map[string]string{},
		},
		{
			name:    "Single variable pattern",
			pattern: "{var}",
			topic:   "value",
			want:    map[string]string{"var": "value"},
		},
		{
			name:    "Pattern with empty segments",
			pattern: "devices/{deviceId}//data",
			topic:   "devices/device123//data",
			want:    map[string]string{"deviceId": "device123"},
		},
		{
			name:    "Topic with empty variable value",
			pattern: "devices/{deviceId}/data",
			topic:   "devices//data",
			want:    map[string]string{"deviceId": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := newTopicParser(tt.pattern)
			if got := parser.extractVariables(tt.topic); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestWildcardTopicIntegration tests the integration between topic patterns and MQTT wildcards
func TestWildcardTopicIntegration(t *testing.T) {
	// Test cases for topic pattern conversion
	tests := []struct {
		name         string
		topicPattern string
		expectedMQTT string
		sampleTopic  string
		expectedVars map[string]string
	}{
		{
			name:         "Simple device pattern",
			topicPattern: "devices/{deviceId}/data",
			expectedMQTT: "devices/+/data",
			sampleTopic:  "devices/device123/data",
			expectedVars: map[string]string{
				"deviceId": "device123",
			},
		},
		{
			name:         "Multiple variables pattern",
			topicPattern: "devices/{deviceId}/sensors/{sensorId}/data",
			expectedMQTT: "devices/+/sensors/+/data",
			sampleTopic:  "devices/device456/sensors/temp1/data",
			expectedVars: map[string]string{
				"deviceId": "device456",
				"sensorId": "temp1",
			},
		},
		{
			name:         "Variable at start",
			topicPattern: "{region}/devices/{deviceId}",
			expectedMQTT: "+/devices/+",
			sampleTopic:  "eu-west/devices/device789",
			expectedVars: map[string]string{
				"region":   "eu-west",
				"deviceId": "device789",
			},
		},
		{
			name:         "Variable at end",
			topicPattern: "users/{userId}/{action}",
			expectedMQTT: "users/+/+",
			sampleTopic:  "users/user123/login",
			expectedVars: map[string]string{
				"userId": "user123",
				"action": "login",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a topic parser
			parser := newTopicParser(tt.topicPattern)

			// Verify MQTT pattern conversion
			mqttPattern := parser.getMQTTSubscriptionTopic()
			assert.Equal(t, tt.expectedMQTT, mqttPattern, "MQTT pattern should match expected")

			// Verify variable extraction
			variables := parser.extractVariables(tt.sampleTopic)
			assert.Equal(t, tt.expectedVars, variables, "Extracted variables should match expected")
		})
	}
}

// TestAdvancedWildcardPatterns tests more complex MQTT wildcard patterns
func TestAdvancedWildcardPatterns(t *testing.T) {
	tests := []struct {
		name         string
		topicPattern string
		mqttPattern  string
		topics       map[string]map[string]string // Map of topic to expected variables
	}{
		{
			name:         "Nested structure",
			topicPattern: "buildings/{buildingId}/floors/{floorId}/rooms/{roomId}/sensors/{sensorId}",
			mqttPattern:  "buildings/+/floors/+/rooms/+/sensors/+",
			topics: map[string]map[string]string{
				"buildings/building1/floors/floor2/rooms/room3/sensors/temp1": {
					"buildingId": "building1",
					"floorId":    "floor2",
					"roomId":     "room3",
					"sensorId":   "temp1",
				},
			},
		},
		{
			name:         "Multiple adjacent variables",
			topicPattern: "{region}/{country}/{city}/weather",
			mqttPattern:  "+/+/+/weather",
			topics: map[string]map[string]string{
				"europe/germany/berlin/weather": {
					"region":  "europe",
					"country": "germany",
					"city":    "berlin",
				},
				"namerica/usa/newyork/weather": {
					"region":  "namerica",
					"country": "usa",
					"city":    "newyork",
				},
			},
		},
		{
			name:         "Variables with special characters",
			topicPattern: "users/{userId}/devices/{deviceId}",
			mqttPattern:  "users/+/devices/+",
			topics: map[string]map[string]string{
				"users/user-123_456/devices/dev.789": {
					"userId":   "user-123_456",
					"deviceId": "dev.789",
				},
			},
		},
		{
			name:         "Single-level wildcard in pattern",
			topicPattern: "devices/{deviceId}/+/{property}",
			mqttPattern:  "devices/+/+/+",
			topics: map[string]map[string]string{
				"devices/device123/temperature/value": {
					"deviceId": "device123",
					"property": "value",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := newTopicParser(tt.topicPattern)

			// Verify MQTT pattern conversion
			assert.Equal(t, tt.mqttPattern, parser.getMQTTSubscriptionTopic(),
				"MQTT pattern should match expected")

			// Test variable extraction for each topic
			for topic, expectedVars := range tt.topics {
				variables := parser.extractVariables(topic)
				assert.Equal(t, expectedVars, variables,
					"Extracted variables for topic %s should match expected", topic)
			}
		})
	}
}

// TestEdgeCasePatterns tests edge cases for topic patterns
func TestEdgeCasePatterns(t *testing.T) {
	tests := []struct {
		name         string
		topicPattern string
		mqttPattern  string
		topic        string
		expectedVars map[string]string
	}{
		{
			name:         "Empty segments",
			topicPattern: "devices/{deviceId}//data",
			mqttPattern:  "devices/+//data",
			topic:        "devices/device123//data",
			expectedVars: map[string]string{
				"deviceId": "device123",
			},
		},
		{
			name:         "Empty variable value",
			topicPattern: "devices/{deviceId}/data",
			mqttPattern:  "devices/+/data",
			topic:        "devices//data",
			expectedVars: map[string]string{
				"deviceId": "",
			},
		},
		{
			name:         "Single variable only",
			topicPattern: "{var}",
			mqttPattern:  "+",
			topic:        "value",
			expectedVars: map[string]string{
				"var": "value",
			},
		},
		{
			name:         "Empty pattern and topic",
			topicPattern: "",
			mqttPattern:  "",
			topic:        "",
			expectedVars: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := newTopicParser(tt.topicPattern)

			// Verify MQTT pattern conversion
			assert.Equal(t, tt.mqttPattern, parser.getMQTTSubscriptionTopic(),
				"MQTT pattern should match expected")

			// Verify variable extraction
			variables := parser.extractVariables(tt.topic)
			assert.Equal(t, tt.expectedVars, variables,
				"Extracted variables should match expected")
		})
	}
}
