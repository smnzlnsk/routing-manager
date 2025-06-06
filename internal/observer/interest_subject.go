package observer

import (
	"sync"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"go.uber.org/zap"
)

// InterestSubject is the concrete implementation of the Subject interface for interest events
type InterestSubject struct {
	observers map[string]domain.Observer
	mutex     sync.RWMutex
	logger    *zap.Logger
}

// NewInterestSubject creates a new instance of InterestSubject
func NewInterestSubject(logger *zap.Logger) *InterestSubject {
	return &InterestSubject{
		observers: make(map[string]domain.Observer),
		logger:    logger,
	}
}

// Register adds an observer to the notification list
func (s *InterestSubject) Register(obs domain.Observer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.observers[obs.GetID()] = obs
	s.logger.Debug("Observer registered")
}

// Deregister removes an observer from the notification list
func (s *InterestSubject) Deregister(obs domain.Observer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.observers, obs.GetID())
	s.logger.Debug("Observer deregistered")
}

// Notify notifies all observers of an event
func (s *InterestSubject) Notify(event domain.InterestEvent) {
	s.mutex.RLock()
	observers := make([]domain.Observer, 0, len(s.observers))
	for _, obs := range s.observers {
		observers = append(observers, obs)
	}
	s.mutex.RUnlock()

	s.logger.Debug("Notifying observers",
		zap.String("eventType", string(event.Type)),
		zap.Int("observerCount", len(observers)))

	for _, obs := range observers {
		go obs.Update(event)
	}
}

// InterestCreated emits an interest created event
func (s *InterestSubject) InterestCreated(interest *domain.Interest) {
	s.Notify(domain.InterestEvent{
		Type:     domain.InterestCreated,
		Interest: interest,
	})
}

// InterestUpdated emits an interest updated event
func (s *InterestSubject) InterestUpdated(interest *domain.Interest) {
	s.Notify(domain.InterestEvent{
		Type:     domain.InterestUpdated,
		Interest: interest,
	})
}

// InterestDeleted emits an interest deleted event
func (s *InterestSubject) InterestDeleted(interest *domain.Interest) {
	s.Notify(domain.InterestEvent{
		Type:     domain.InterestDeleted,
		Interest: interest,
	})
}
