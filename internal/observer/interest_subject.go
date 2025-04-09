package observer

import (
	"sync"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/domain/observer"
	"go.uber.org/zap"
)

// InterestSubject is the concrete implementation of the Subject interface for interest events
type InterestSubject struct {
	observers []observer.Observer
	mutex     sync.RWMutex
	logger    *zap.Logger
}

// NewInterestSubject creates a new instance of InterestSubject
func NewInterestSubject(logger *zap.Logger) *InterestSubject {
	return &InterestSubject{
		observers: make([]observer.Observer, 0),
		logger:    logger,
	}
}

// Register adds an observer to the notification list
func (s *InterestSubject) Register(obs observer.Observer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.observers = append(s.observers, obs)
	s.logger.Debug("Observer registered")
}

// Deregister removes an observer from the notification list
func (s *InterestSubject) Deregister(obs observer.Observer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, o := range s.observers {
		if o == obs {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			s.logger.Debug("Observer deregistered")
			return
		}
	}
}

// Notify notifies all observers of an event
func (s *InterestSubject) Notify(event observer.InterestEvent) {
	s.mutex.RLock()
	observers := make([]observer.Observer, len(s.observers))
	copy(observers, s.observers)
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
	s.Notify(observer.InterestEvent{
		Type:     observer.InterestCreated,
		Interest: interest,
	})
}

// InterestUpdated emits an interest updated event
func (s *InterestSubject) InterestUpdated(interest *domain.Interest) {
	s.Notify(observer.InterestEvent{
		Type:     observer.InterestUpdated,
		Interest: interest,
	})
}

// InterestDeleted emits an interest deleted event
func (s *InterestSubject) InterestDeleted(interest *domain.Interest) {
	s.Notify(observer.InterestEvent{
		Type:     observer.InterestDeleted,
		Interest: interest,
	})
}
