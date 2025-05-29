package repositories

import (
	"context"
	"time"

	"ui_api.go/internal/models"
)

type LogFilter struct {
	IP     string
	Method string
	Path   string
	Status *int
	From   *time.Time
	To     *time.Time
}

type LogRepository interface {
	FindByFilter(ctx context.Context, filter LogFilter) ([]models.LogEntry, error)
	GetFilteredLogs(ctx context.Context, filter LogFilter, limit int, offset int) ([]models.LogEntry, int64, error)
	CountByStatus(ctx context.Context, status int) (int64, error)
	GetErrorStats(ctx context.Context, from, to time.Time) (map[int]int64, error)
	GetTrafficStats(ctx context.Context, from, to time.Time) (total, unique int64, err error)
	GetLatencyStats(ctx context.Context, from, to time.Time) (avg float64, max int64, err error)
}
