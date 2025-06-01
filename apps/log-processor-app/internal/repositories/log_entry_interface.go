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
}
