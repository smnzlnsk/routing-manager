package service

import (
	"context"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.uber.org/zap"
)

type RoutingService interface {
	HandleRoutingChange(ctx context.Context, routingChange *domain.RoutingChange) error
	GetRouting(ctx context.Context, appName string) (*domain.Job, error)
}

type routingService struct {
	repo   repository.RoutingRepository
	logger *zap.Logger
}

func NewRoutingService(repo repository.RoutingRepository, logger *zap.Logger) RoutingService {
	return &routingService{
		repo:   repo,
		logger: logger,
	}
}

func (s *routingService) HandleRoutingChange(ctx context.Context, routingChange *domain.RoutingChange) error {
	s.logger.Info("Handling routing change", zap.Any("routingChange", routingChange))
	return nil
}

func (s *routingService) GetRouting(ctx context.Context, appName string) (*domain.Job, error) {
	s.logger.Info("Getting routing", zap.String("appName", appName))
	return nil, nil
}
