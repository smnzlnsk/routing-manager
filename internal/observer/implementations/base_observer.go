package implementations

import (
	"github.com/smnzlnsk/routing-manager/internal/domain"
	"go.uber.org/zap"
)

// BaseObserver provides common functionality for observers
type BaseObserver struct {
	name   string
	logger *zap.Logger
}

var _ domain.Observer = &BaseObserver{}

// NewBaseObserver creates a new BaseObserver
func NewBaseObserver(name string, logger *zap.Logger) *BaseObserver {
	return &BaseObserver{
		name:   name,
		logger: logger,
	}
}

// Update is called when an event occurs
// This is a placeholder implementation that should be overridden by concrete observers
func (o *BaseObserver) Update(event domain.InterestEvent) {
	o.logger.Debug("Observer received event",
		zap.String("observer", o.name),
		zap.String("eventType", string(event.Type)))
}

// GetID returns the ID of the observer's interest
func (o *BaseObserver) GetID() string {
	return o.name
}
