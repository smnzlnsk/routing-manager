package repository

import (
	"context"

	"github.com/smnzlnsk/routing-manager/internal/domain"
)

type JobRepository interface {
	GetByJobName(ctx context.Context, jobName string) (*domain.Job, error)
}
