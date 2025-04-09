package sql

import (
	"context"
	"database/sql"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.uber.org/zap"
)

type interestRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewInterestRepository(db *sql.DB, logger *zap.Logger) repository.InterestRepository {
	return &interestRepository{
		db:     db,
		logger: logger,
	}
}

func (r *interestRepository) Create(ctx context.Context, interest *domain.Interest) error {
	return nil
}

func (r *interestRepository) GetByAppName(ctx context.Context, appName string) (*domain.Interest, error) {
	return nil, nil
}

func (r *interestRepository) GetByServiceIp(ctx context.Context, serviceIp string) (*domain.Interest, error) {
	return nil, nil
}

func (r *interestRepository) Update(ctx context.Context, interest *domain.Interest) (*domain.Interest, error) {
	return nil, nil
}

func (r *interestRepository) DeleteByAppName(ctx context.Context, appName string) error {
	return nil
}

func (r *interestRepository) DeleteByServiceIp(ctx context.Context, serviceIp string) error {
	return nil
}

func (r *interestRepository) List(ctx context.Context) ([]*domain.Interest, error) {
	return nil, nil
}
