package main

import (
	"sync"
	"time"

	"generator.go/internal/config"
	"generator.go/internal/service"
)

func main() {

	config := config.InitConfig()

	var wg sync.WaitGroup

	ticker := time.NewTicker(time.Second / time.Duration(config.RPS))
	defer ticker.Stop()

	logGenerationService := service.NewLogGenerationService(config)
	sendLogService := service.NewSendLogService(logGenerationService)

	for range ticker.C {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendLogService.SendLogs(config.TargetApiUrl)
		}()
	}
}
