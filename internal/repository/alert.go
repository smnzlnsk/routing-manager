package repository

import (
	"context"

	"github.com/smnzlnsk/routing-manager/internal/domain"
)

type AlertRepository interface {
	Create(ctx context.Context, alert *domain.Alert) error
	GetByAppName(ctx context.Context, appName string) (*domain.Alert, error)
}
