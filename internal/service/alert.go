package service

import (
	"context"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.uber.org/zap"
)

type AlertService interface {
	HandleAlert(ctx context.Context, alert *domain.Alert) error
}

type alertService struct {
	repo   repository.AlertRepository
	logger *zap.Logger
}

func NewAlertService(repo repository.AlertRepository, logger *zap.Logger) AlertService {
	return &alertService{
		repo:   repo,
		logger: logger,
	}
}

func (s *alertService) HandleAlert(ctx context.Context, alert *domain.Alert) error {
	s.logger.Info("Handling alert", zap.Any("alert", alert))

	return nil
}
