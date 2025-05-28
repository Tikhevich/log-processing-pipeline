package repositories

import (
	"context"
	"time"

	"worker.go/internal/models"
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
	Create(ctx context.Context, log *models.LogEntry) error
	GetByID(ctx context.Context, id int64) (*models.LogEntry, error)
	FindByFilter(ctx context.Context, filter LogFilter) ([]models.LogEntry, error)
	DeleteOldLogs(ctx context.Context, olderThan time.Time) error
}
