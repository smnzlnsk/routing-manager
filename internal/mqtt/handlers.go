package mqtt

import (
	"encoding/json"
	"fmt"

	"github.com/smnzlnsk/routing-manager/internal/logger"
	"github.com/smnzlnsk/routing-manager/internal/models"
)

// HandleRoutingManagerMessage handles messages from the routing manager topic
func HandleRoutingManagerMessage(topic string, payload []byte) {
	logger.Info("Received message from routing manager topic", "topic", topic, "payload", string(payload))
}

// HandleTableQueryMessage handles messages from table query topics with variables
// Topic pattern: tablequery/{requesterId}/{routingPolicy}
func HandleTableQueryMessage(topic string, variables map[string]string, payload []byte) {
	requesterId := variables["requesterId"]
	routingPolicy := variables["routingPolicy"]

	logger.Infof("Received table query message:\n topic: %s\n requesterId: %s\n routingPolicy: %s\n payload: %s",
		topic,
		requesterId,
		routingPolicy,
		string(payload))

	request := models.TableQueryMessageRequest{
		RequesterId:   requesterId,
		RoutingPolicy: routingPolicy,
		Payload:       payload,
	}

	response, err := handleTableQueryMessage(request)
	if err != nil {
		logger.Errorf("Failed to handle table query message: %v", err)
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("Failed to marshal table query message response: %v", err)
	}

	responseTopic := fmt.Sprintf("result/%s", topic)
	err = Instance().Publish(responseTopic, responseBytes)
	if err != nil {
		logger.Errorf("Failed to publish table query message response: %v", err)
	}
}

// DebugHandler is a generic handler for any topic with wildcards
func DebugHandler(topic string, variables map[string]string, payload []byte) {
	logger.Infof("Received message on wildcard topic:\n topic: %s\n variables: %v\n payload: %s",
		topic,
		variables,
		string(payload))
}
