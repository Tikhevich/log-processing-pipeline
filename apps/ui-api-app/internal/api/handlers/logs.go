package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"ui_api.go/internal/models"
	"ui_api.go/internal/repositories"
	"ui_api.go/internal/service"
)

type LogsHandler struct {
	logsService *service.LogsService
}

type LogsResponse struct {
	Data []models.LogEntry `json:"data"`
	Meta struct {
		Total  int64 `json:"total"`
		Limit  int   `json:"limit"`
		Offset int   `json:"offset"`
	} `json:"meta"`
}

func NewLogsHandler(logsService *service.LogsService) *LogsHandler {
	return &LogsHandler{logsService: logsService}
}

func (logsHandler *LogsHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/logs", logsHandler.GetLogs)
}

func (h *LogsHandler) GetLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filter := getPreparedFilterBasedonContext(c)
	ctx := context.Background()

	logs, total, err := h.logsService.GetLogsByParams(ctx, filter, limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := LogsResponse{
		Data: logs,
	}
	response.Meta.Total = total
	response.Meta.Limit = limit
	response.Meta.Offset = offset

	c.JSON(http.StatusOK, response)
}

func getPreparedFilterBasedonContext(c *gin.Context) repositories.LogFilter {
	filter := repositories.LogFilter{
		IP:     c.Query("ip"),
		Method: c.Query("method"),
		Path:   c.Query("path"),
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err == nil {
			filter.Status = &status
		}
	}

	if fromStr := c.Query("from"); fromStr != "" {
		if from, err := time.Parse(time.RFC3339, fromStr); err == nil {
			filter.From = &from
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if to, err := time.Parse(time.RFC3339, toStr); err == nil {
			filter.To = &to
		}
	}

	return filter
}
