package implementations

import (
	"github.com/smnzlnsk/routing-manager/internal/domain/observer"
	"go.uber.org/zap"
)

// BaseObserver provides common functionality for observers
type BaseObserver struct {
	name   string
	logger *zap.Logger
}

// NewBaseObserver creates a new BaseObserver
func NewBaseObserver(name string, logger *zap.Logger) *BaseObserver {
	return &BaseObserver{
		name:   name,
		logger: logger,
	}
}

// Update is called when an event occurs
// This is a placeholder implementation that should be overridden by concrete observers
func (o *BaseObserver) Update(event observer.InterestEvent) {
	o.logger.Debug("Observer received event",
		zap.String("observer", o.name),
		zap.String("eventType", string(event.Type)))
}
