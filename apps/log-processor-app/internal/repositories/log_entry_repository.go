package repositories

import (
	"context"

	"gorm.io/gorm"
	"worker.go/internal/models"
)

type logRepository struct {
	db *gorm.DB
}

func GetLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

// Create implements LogRepository.
func (r *logRepository) Create(ctx context.Context, log *models.LogEntry) error {
	return r.db.WithContext(ctx).Create(log).Error
}
