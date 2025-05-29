package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"ui_api.go/internal/service"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (statsHandler *StatsHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/stats", statsHandler.GetStats)
}

func (statsHandler *StatsHandler) GetStats(c *gin.Context) {
	ctx := context.Background()

	rangeType := c.DefaultQuery("range", "hour")
	statsType := c.DefaultQuery("type", "errors")

	stats, err := statsHandler.statsService.GetStatsByStatsTypeAndRange(ctx, statsType, rangeType)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stats type"})
	}

	c.JSON(http.StatusOK, stats)
}
