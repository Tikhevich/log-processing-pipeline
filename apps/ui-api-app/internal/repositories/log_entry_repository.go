package repositories

import (
	"context"
	"net/http"
	"time"

	"gorm.io/gorm"
	"ui_api.go/internal/models"
)

type logRepository struct {
	db *gorm.DB
}

// GetErrorStats implements LogRepository.
func (r *logRepository) GetErrorStats(from time.Time, to time.Time) (map[int]int64, error) {
	var results []struct {
		Status int
		Count  int64
	}

	err := r.db.Model(&models.LogEntry{}).
		Select("status, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", from, to).
		Where("status <> ?", http.StatusOK).
		Group("status").
		Scan(&results).
		Error

	stats := make(map[int]int64)
	for _, res := range results {
		stats[res.Status] = res.Count
	}

	return stats, err
}

// GetLatencyStats implements LogRepository.
func (r *logRepository) GetLatencyStats(from time.Time, to time.Time) (avg float64, max int64, err error) {
	var result struct {
		Avg float64
		Max int64
	}

	err = r.db.Model(&models.LogEntry{}).
		Select("AVG(latency_ms) as avg, MAX(latency_ms) as max").
		Where("timestamp BETWEEN ? AND ?", from, to).
		Scan(&result).
		Error

	return result.Avg, result.Max, err
}

// GetTrafficStats implements LogRepository.
func (r *logRepository) GetTrafficStats(from time.Time, to time.Time) (total int64, unique int64, err error) {
	var result struct {
		UniqueIps     int64
		TotalRequests int64
	}

	err = r.db.Model(&models.LogEntry{}).
		Select("COUNT(DISTINCT ip) as unique_ips, COUNT(*) as total_requests").
		Where("timestamp BETWEEN ? AND ?", from, to).
		Scan(&result).
		Error

	return result.TotalRequests, result.UniqueIps, err
}

// GetFilteredLogs implements LogRepository.
func (r *logRepository) GetFilteredLogs(ctx context.Context, filter LogFilter, limit int, offset int) ([]models.LogEntry, int64, error) {
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

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&logs).Error
	if err != nil {
		return nil, total, err
	}

	return logs, total, nil
}

// CountByStatus implements LogRepository.
func (r *logRepository) CountByStatus(ctx context.Context, status int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.LogEntry{}).
		Where("status = ?", status).
		Count(&count).
		Error
	return count, err
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
