package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ui_api.go/internal/api/handlers"
	"ui_api.go/internal/repositories"
	"ui_api.go/internal/service"
)

func main() {

	db, err := gorm.Open(mysql.Open("app_user:app_pass@tcp(localhost:3306)/log_db?charset=utf8mb4&parseTime=True"))
	if err != nil {
		fmt.Println(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	repo := repositories.GetLogRepository(db)

	statsService := service.NewStatsService(repo, redisClient)
	statsHandler := handlers.NewStatsService(statsService)

	router := gin.Default()
	statsHandler.RegisterRoutes(router)

	router.GET("/logs", handlers.GetLogs)
	router.Run(":8081")
}
