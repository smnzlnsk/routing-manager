package repository

import (
	"context"

	"github.com/smnzlnsk/routing-manager/internal/domain"
)

type RoutingRepository interface {
	GetRouting(ctx context.Context, appName string) (*domain.JobRouting, error)
	UpdateRouting(ctx context.Context, routing *domain.JobRouting) error
}
