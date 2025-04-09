package sql

import (
	"database/sql"

	"github.com/smnzlnsk/routing-manager/internal/repository"
	"go.uber.org/zap"
)

func New(db *sql.DB, logger *zap.Logger) *repository.Repositories {
	return &repository.Repositories{
		InterestRepository: NewInterestRepository(db, logger),
		// Initialize other repositories here with their dependencies
	}
}
