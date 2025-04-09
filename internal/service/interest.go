package service

import (
	"context"
	"errors"
	"time"

	"github.com/smnzlnsk/routing-manager/internal/domain"
	"github.com/smnzlnsk/routing-manager/internal/domain/observer"
	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.uber.org/zap"
)

type InterestService interface {
	Create(ctx context.Context, interest *domain.Interest) (*domain.Interest, error)
	GetByAppName(ctx context.Context, appName string) (*domain.Interest, error)
	GetByServiceIp(ctx context.Context, serviceIp string) (*domain.Interest, error)
	Update(ctx context.Context, interest *domain.Interest) (*domain.Interest, error)
	DeleteByAppName(ctx context.Context, appName string) error
	DeleteByServiceIp(ctx context.Context, serviceIp string) error
	List(ctx context.Context) ([]*domain.Interest, error)
}

type interestService struct {
	repo    repository.InterestRepository
	logger  *zap.Logger
	subject observer.Subject
}

func NewInterestService(repo repository.InterestRepository, subject observer.Subject, logger *zap.Logger) InterestService {
	return &interestService{
		repo:    repo,
		logger:  logger,
		subject: subject,
	}
}

func (s *interestService) Create(ctx context.Context, interest *domain.Interest) (*domain.Interest, error) {
	s.logger.Info("Creating interest", zap.Any("interest", interest))

	// Check if the interest already exists
	existingInterest, err := s.repo.GetByAppName(ctx, interest.AppName)
	if err != nil {
		s.logger.Error("Error in create interest", zap.Error(err))
		var domainErr *domain.Error
		if !errors.As(err, &domainErr) || domainErr.Code != domain.CodeNotFound {
			s.logger.Error("Returning error in create interest", zap.Error(err))
			return nil, err
		}
	}

	if existingInterest != nil {
		s.logger.Error("Interest already exists", zap.Any("interest", existingInterest))
		return nil, domain.ErrInterestAlreadyExists
	}

	now := time.Now()
	i := &domain.Interest{
		AppName:   interest.AppName,
		ServiceIp: interest.ServiceIp,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.logger.Info("Creating interest in repo", zap.Any("interest", i))

	if err := s.repo.Create(ctx, i); err != nil {
		s.logger.Error("Error in repo create interest", zap.Error(err))
		return nil, err
	}

	// Notify observers about the created interest
	if s.subject != nil {
		s.subject.Notify(observer.InterestEvent{
			Type:     observer.InterestCreated,
			Interest: i,
		})
	}

	return i, nil
}

func (s *interestService) GetByAppName(ctx context.Context, appName string) (*domain.Interest, error) {
	s.logger.Debug("Getting interest by app name", zap.String("appName", appName))
	return s.repo.GetByAppName(ctx, appName)
}

func (s *interestService) GetByServiceIp(ctx context.Context, serviceIp string) (*domain.Interest, error) {
	s.logger.Debug("Getting interest by service IP", zap.String("serviceIp", serviceIp))
	return s.repo.GetByServiceIp(ctx, serviceIp)
}

func (s *interestService) Update(ctx context.Context, interest *domain.Interest) (*domain.Interest, error) {
	s.logger.Debug("Updating interest", zap.Any("interest", interest))
	updatedInterest, err := s.repo.Update(ctx, interest)
	if err != nil {
		return nil, err
	}
	return updatedInterest, nil
}

func (s *interestService) DeleteByAppName(ctx context.Context, appName string) error {
	s.logger.Debug("Deleting interest by app name", zap.String("appName", appName))
	return s.repo.DeleteByAppName(ctx, appName)
}

func (s *interestService) DeleteByServiceIp(ctx context.Context, serviceIp string) error {
	s.logger.Debug("Deleting interest by service IP", zap.String("serviceIp", serviceIp))
	return s.repo.DeleteByServiceIp(ctx, serviceIp)
}

func (s *interestService) List(ctx context.Context) ([]*domain.Interest, error) {
	s.logger.Debug("Listing interests")
	return s.repo.List(ctx)
}
