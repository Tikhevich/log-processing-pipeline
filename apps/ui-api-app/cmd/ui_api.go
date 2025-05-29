package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ui_api.go/internal/api/handlers"
	"ui_api.go/internal/config"
	"ui_api.go/internal/repositories"
	"ui_api.go/internal/service"
)

func main() {

	config := config.Load()

	db, err := gorm.Open(mysql.Open(config.DBDSN))
	if err != nil {
		fmt.Println(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		DB:   0,
	})

	repo := repositories.GetLogRepository(db)

	statsService := service.NewStatsService(repo, redisClient)
	statsHandler := handlers.NewStatsHandler(statsService)

	logService := service.NewLogsService(repo)
	logsHandler := handlers.NewLogsHandler(logService)

	router := gin.Default()
	statsHandler.RegisterRoutes(router)
	logsHandler.RegisterRoutes(router)

	router.Run(":8081")
}
