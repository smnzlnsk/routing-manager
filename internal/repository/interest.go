package repository

import (
	"context"

	"github.com/smnzlnsk/routing-manager/internal/domain"
)

type InterestRepository interface {
	Create(ctx context.Context, interest *domain.Interest) error
	GetByAppName(ctx context.Context, appName string) (*domain.Interest, error)
	GetByServiceIp(ctx context.Context, serviceIp string) (*domain.Interest, error)
	Update(ctx context.Context, interest *domain.Interest) (*domain.Interest, error)
	DeleteByAppName(ctx context.Context, appName string) error
	DeleteByServiceIp(ctx context.Context, serviceIp string) error
	List(ctx context.Context) ([]*domain.Interest, error)
}
