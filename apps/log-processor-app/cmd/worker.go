package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"worker.go/internal/config"
	"worker.go/internal/consumer"
	"worker.go/internal/models"
	"worker.go/internal/repositories"
	"worker.go/internal/service"
)

func main() {
	workerConfig := config.Load()

	db, err := gorm.Open(mysql.Open(workerConfig.DBDSN))
	if err != nil {
		panic(fmt.Errorf("Couldn't connect to database: %w", err))
	}

	err = db.AutoMigrate(&models.LogEntry{})
	if err != nil {
		panic(fmt.Errorf("Failed to migrate database: %w", err))
	}

	fmt.Println("Waiting for messages...")

	logRepo := repositories.GetLogRepository(db)

	brokerMessageProcessingService := service.NewBrokerMessageProcessingService(logRepo)

	kafkaConsumer := consumer.New(workerConfig)

	for {
		message, err := kafkaConsumer.Read()
		if err != nil {
			break
		}

		stats := brokerMessageProcessingService.ProcessBrokerMessage(message.Value)

		fmt.Printf("Processed log count: %v \n", stats.Processed)
		fmt.Printf("Count of logs except 200 ok status : %v \n", stats.Added)
	}

	err = kafkaConsumer.Close()
	if err != nil {
		panic(fmt.Errorf("Failed to close reader: %w", err))
	}
}
