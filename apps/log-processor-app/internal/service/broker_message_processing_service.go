package service

import (
	"context"
	"fmt"
	"net/http"

	"worker.go/internal/parser"
	"worker.go/internal/repositories"
)

type Stats struct {
	Processed int
	Added     int
}

type BrokerMessageProcessingService struct {
	repo repositories.LogRepository
}

func NewBrokerMessageProcessingService(repo repositories.LogRepository) *BrokerMessageProcessingService {
	return &BrokerMessageProcessingService{repo: repo}
}

func (b *BrokerMessageProcessingService) ProcessBrokerMessage(messageValue []byte) Stats {
	newLogs, err := parser.ParseLogs(messageValue)
	if err != nil {
		fmt.Printf("Failed to parse log: %v", err)
	}

	ctx := context.Background()

	var countOfAddedLogs int
	for _, newLog := range newLogs {
		// Skip 200 ok status logs
		if newLog.Status == http.StatusOK {
			continue
		}

		err = b.repo.Create(ctx, &newLog)
		if err != nil {
			fmt.Printf("Failed to create log: %v", err)
		}
		countOfAddedLogs++
	}

	return Stats{Processed: len(newLogs), Added: countOfAddedLogs}
}
