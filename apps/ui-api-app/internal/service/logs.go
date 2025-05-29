package service

import (
	"context"

	"ui_api.go/internal/models"
	"ui_api.go/internal/repositories"
)

type LogsService struct {
	repo repositories.LogRepository
}

func NewLogsService(repo repositories.LogRepository) *LogsService {
	return &LogsService{repo: repo}
}

func (s *LogsService) GetLogsByParams(ctx context.Context, filter repositories.LogFilter, limit int, offset int) ([]models.LogEntry, int64, error) {
	logs, total, err := s.repo.GetFilteredLogs(ctx, filter, limit, offset)

	return logs, total, err
}
