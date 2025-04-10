package service

import (
	"context"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.uber.org/zap"
)

type JobService interface {
	GetByJobName(ctx context.Context, jobName string) (*domain.Job, error)
}

type jobService struct {
	repo   repository.JobRepository
	logger *zap.Logger
}

func NewJobService(repo repository.JobRepository, logger *zap.Logger) JobService {
	return &jobService{
		repo:   repo,
		logger: logger,
	}
}

func (s *jobService) GetByJobName(ctx context.Context, jobName string) (*domain.Job, error) {
	s.logger.Info("Getting job by job name", zap.String("jobName", jobName))
	return s.repo.GetByJobName(ctx, jobName)
}
