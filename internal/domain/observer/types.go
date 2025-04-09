package observer

import "github.com/smnzlnsk/routing-manager/internal/domain"

// EventType defines the type of event
type EventType string

// Event types
const (
	InterestCreated EventType = "INTEREST_CREATED"
	InterestUpdated EventType = "INTEREST_UPDATED"
	InterestDeleted EventType = "INTEREST_DELETED"
)

// InterestEvent represents an event related to an interest
type InterestEvent struct {
	Type     EventType
	Interest *domain.Interest
}

// Observer defines the interface for objects that want to be notified of events
type Observer interface {
	// Update is called when an event occurs
	Update(event InterestEvent)
}

// Subject defines the interface for objects that maintain observers
type Subject interface {
	// Register adds an observer to the notification list
	Register(observer Observer)

	// Deregister removes an observer from the notification list
	Deregister(observer Observer)

	// Notify notifies all observers of an event
	Notify(event InterestEvent)
}
