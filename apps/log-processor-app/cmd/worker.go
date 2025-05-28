package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"worker.go/internal/config"
	"worker.go/internal/consumer"
	"worker.go/internal/models"
	"worker.go/internal/parser"
	"worker.go/internal/repositories"
)

func main() {
	workerConfig := config.Load()

	db, err := gorm.Open(mysql.Open(workerConfig.DBDSN))
	if err != nil {
		fmt.Println(err)
	}

	err = db.AutoMigrate(&models.LogEntry{})
	if err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	logRepo := repositories.GetLogRepository(db)

	kafkaConsumer := consumer.New(workerConfig)

	for {
		message, err := kafkaConsumer.Read()
		if err != nil {
			break
		}

		newLogs, err := parser.ParseLogs(message.Value)
		if err != nil {
			log.Printf("Failed to parse log: %v", err)
		}

		ctx := context.Background()

		var countOfAddedLogs int
		for _, newLog := range newLogs {
			// Skip 200 ok status logs
			if newLog.Status == http.StatusOK {
				continue
			}

			err = logRepo.Create(ctx, &newLog)
			if err != nil {
				log.Printf("Failed to create log: %v", err)
			}
			countOfAddedLogs++
		}
		fmt.Printf("Length of new logs : %v \n", len(newLogs))
		fmt.Printf("Actualy added logs : %v \n", countOfAddedLogs)
	}

	err = kafkaConsumer.Close()
	if err != nil {
		log.Fatal("failed to close reader:", err)
	}

}
