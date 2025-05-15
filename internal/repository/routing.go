package repository

import (
	"context"

	"github.com/smnzlnsk/routing-manager/internal/domain"
)

type RoutingRepository interface {
	GetRouting(ctx context.Context, appName string) (*domain.Job, error)
	UpdateRouting(ctx context.Context, routing *domain.Job) error
}
