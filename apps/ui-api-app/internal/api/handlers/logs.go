package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ui_api.go/internal/models"
	"ui_api.go/internal/repositories"
)

type LogsResponse struct {
	Data []models.LogEntry `json:"data"`
	Meta struct {
		Total  int64 `json:"total"`
		Limit  int   `json:"limit"`
		Offset int   `json:"offset"`
	} `json:"meta"`
}

func GetLogs(c *gin.Context) {

	db, err := gorm.Open(mysql.Open("app_user:app_pass@tcp(localhost:3306)/log_db?charset=utf8mb4&parseTime=True"))
	if err != nil {
		fmt.Println(err)
	}

	logRepo := repositories.GetLogRepository(db)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

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

	ctx := context.Background()

	logs, total, err := logRepo.GetFilteredLogs(ctx, filter, limit, offset)
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
