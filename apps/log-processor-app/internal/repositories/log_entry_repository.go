package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"
	"worker.go/internal/models"
)

type logRepository struct {
	db *gorm.DB
}

// DeleteOldLogs implements LogRepository.
func (r *logRepository) DeleteOldLogs(ctx context.Context, olderThan time.Time) error {
	return r.db.WithContext(ctx).
		Where("timestamp < ?", olderThan).
		Delete(&models.LogEntry{}).
		Error
}

// FindByFilter implements LogRepository.
func (r *logRepository) FindByFilter(ctx context.Context, filter LogFilter) ([]models.LogEntry, error) {
	var logs []models.LogEntry
	query := r.db.WithContext(ctx).Model(&models.LogEntry{})

	if filter.IP != "" {
		query = query.Where("ip = ?", filter.IP)
	}
	if filter.Method != "" {
		query = query.Where("method = ?", filter.Method)
	}
	if filter.Path != "" {
		query = query.Where("path LIKE ?", "%"+filter.Path+"%")
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.From != nil {
		query = query.Where("timestamp >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("timestamp <= ?", *filter.To)
	}

	err := query.Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetByID implements LogRepository.
func (r *logRepository) GetByID(ctx context.Context, id int64) (*models.LogEntry, error) {
	var log models.LogEntry
	err := r.db.WithContext(ctx).First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func GetLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

// Create implements LogRepository.
func (r *logRepository) Create(ctx context.Context, log *models.LogEntry) error {
	return r.db.WithContext(ctx).Create(log).Error
}
