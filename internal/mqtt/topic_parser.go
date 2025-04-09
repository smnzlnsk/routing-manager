package mqtt

import (
	"regexp"
	"strings"
)

// topicParser handles parsing MQTT topics with named wildcards
type topicParser struct {
	pattern          string
	mqttPattern      string
	variableNames    []string
	variablePatterns []string
}

// newTopicParser creates a new topic parser for the given pattern
// Pattern format: "devices/{deviceId}/sensors/{sensorId}/data"
// MQTT subscription format: "devices/+/sensors/+/data"
func newTopicParser(pattern string) *topicParser {
	parser := &topicParser{
		pattern:          pattern,
		variableNames:    make([]string, 0),
		variablePatterns: make([]string, 0),
	}

	// Extract variable names and convert to MQTT pattern
	mqttPattern := pattern

	// Find all variables in the format {variableName}
	re := regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`)
	matches := re.FindAllStringSubmatch(pattern, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			fullMatch := match[0] // e.g., {deviceId}
			varName := match[1]   // e.g., deviceId

			// Add to variable names
			parser.variableNames = append(parser.variableNames, varName)

			// Replace in MQTT pattern with + wildcard
			mqttPattern = strings.Replace(mqttPattern, fullMatch, "+", 1)

			// Create a pattern to extract this variable later
			// Replace the current variable with a capture group and all other variables with wildcards
			varPattern := pattern
			for _, otherMatch := range matches {
				if otherMatch[0] != fullMatch {
					varPattern = strings.Replace(varPattern, otherMatch[0], "([^/]+)", 1)
				}
			}
			varPattern = strings.Replace(varPattern, fullMatch, "([^/]+)", 1)
			// Convert to regex pattern
			varPattern = strings.Replace(varPattern, "/", "\\/", -1)
			varPattern = "^" + varPattern + "$"

			parser.variablePatterns = append(parser.variablePatterns, varPattern)
		}
	}

	parser.mqttPattern = mqttPattern
	return parser
}

// getMQTTSubscriptionTopic returns the MQTT subscription topic with + wildcards
func (p *topicParser) getMQTTSubscriptionTopic() string {
	return p.mqttPattern
}

// extractVariables extracts variables from a topic based on the pattern
func (p *topicParser) extractVariables(topic string) map[string]string {
	variables := make(map[string]string)

	// Simple case: if there are no variables, return empty map
	if len(p.variableNames) == 0 {
		return variables
	}

	// Split the topic and pattern into segments
	topicSegments := strings.Split(topic, "/")
	patternSegments := strings.Split(p.pattern, "/")

	// If lengths don't match, can't extract variables
	if len(topicSegments) != len(patternSegments) {
		return variables
	}

	// Extract variables by matching segments
	for i, patternSegment := range patternSegments {
		if strings.HasPrefix(patternSegment, "{") && strings.HasSuffix(patternSegment, "}") {
			// This is a variable segment
			varName := patternSegment[1 : len(patternSegment)-1]
			variables[varName] = topicSegments[i]
		}
	}

	return variables
}
