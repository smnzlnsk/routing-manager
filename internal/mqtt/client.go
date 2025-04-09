package mqtt

import (
	"fmt"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/smnzlnsk/routing-manager/config"
	"github.com/smnzlnsk/routing-manager/internal/logger"
)

// MessageHandler is a function that processes MQTT messages
type MessageHandler func(topic string, payload []byte)

// WildcardMessageHandler is a function that processes MQTT messages with topic variables
type WildcardMessageHandler func(topic string, variables map[string]string, payload []byte)

// Client defines the interface for MQTT communication
type Client interface {
	// Connect establishes a connection to the MQTT broker
	Connect() error

	// Disconnect closes the connection to the MQTT broker
	Disconnect()

	// Subscribe subscribes to a topic and registers a handler for messages
	Subscribe(topic string, handler MessageHandler) error

	// SubscribeWithWildcards subscribes to a topic with wildcards and extracts variables from the topic
	SubscribeWithWildcards(topicPattern string, handler WildcardMessageHandler) error

	// Publish publishes a message to a topic
	Publish(topic string, payload []byte) error

	// IsConnected returns true if connected to the broker
	IsConnected() bool
}

// client implements the MQTTBroker interface
type client struct {
	config config.MQTTConfig
	client mqtt.Client
}

// Singleton instance and mutex for thread safety
var (
	instance *client
	once     sync.Once
	mu       sync.Mutex
)

// Instance returns the singleton instance of the MQTT client
// If the instance is not initialized, it will panic
func Instance() *client {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		logger.Fatal("MQTT client not initialized. Make sure to call InitInstance first.")
	}

	return instance
}

// InitInstance initializes the singleton instance with the provided config
// This should be called once at application startup
func InitInstance(config config.MQTTConfig) {
	once.Do(func() {
		mu.Lock()
		defer mu.Unlock()

		instance = &client{
			config: config,
		}
	})
}

// Connect establishes a connection to the MQTT broker
func (c *client) Connect() error {
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", c.config.Host, c.config.Port)).
		SetClientID(c.config.ClientID).
		SetUsername(c.config.Username).
		SetPassword(c.config.Password).
		SetCleanSession(c.config.CleanSession).
		SetConnectTimeout(c.config.ConnectTimeout).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(client mqtt.Client) {
			logger.Info("Connected to MQTT broker")
		}).
		SetConnectionLostHandler(func(client mqtt.Client, err error) {
			logger.Errorf("Connection to MQTT broker lost: %v", err)
		})

	c.client = mqtt.NewClient(opts)
	token := c.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	return nil
}

// Disconnect closes the connection to the MQTT broker
func (c *client) Disconnect() {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250) // 250ms timeout
	}
}

// Subscribe subscribes to a topic and registers a handler for messages
func (c *client) Subscribe(topic string, handler MessageHandler) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected to MQTT broker")
	}

	token := c.client.Subscribe(topic, c.config.QoS, func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Topic(), msg.Payload())
	})

	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}

	logger.Infof("Subscribed to topic: %s", topic)
	return nil
}

// SubscribeWithWildcards subscribes to a topic with wildcards and extracts variables from the topic
func (c *client) SubscribeWithWildcards(topicPattern string, handler WildcardMessageHandler) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected to MQTT broker")
	}

	// Create a topic parser for this pattern
	parser := newTopicParser(topicPattern)

	// Subscribe to the MQTT topic with wildcards
	token := c.client.Subscribe(parser.getMQTTSubscriptionTopic(), c.config.QoS, func(client mqtt.Client, msg mqtt.Message) {
		// Extract variables from the topic
		variables := parser.extractVariables(msg.Topic())

		// Call the handler with the extracted variables
		handler(msg.Topic(), variables, msg.Payload())
	})

	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic pattern %s: %w", topicPattern, token.Error())
	}

	logger.Infof("Subscribed to topic pattern: %s", topicPattern)
	return nil
}

// Publish publishes a message to a topic
func (c *client) Publish(topic string, payload []byte) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected to MQTT broker")
	}

	token := c.client.Publish(topic, c.config.QoS, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish to topic %s: %w", topic, token.Error())
	}

	return nil
}

// IsConnected returns true if connected to the broker
func (c *client) IsConnected() bool {
	return c.client != nil && c.client.IsConnected()
}
