package service

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"generator.go/internal/config"
)

type LogGenerationService struct {
	config config.Config
}

func NewLogGenerationService(config config.Config) *LogGenerationService {
	return &LogGenerationService{config: config}
}

func (s LogGenerationService) GenerateLogs() []string {
	var logs = make([]string, s.config.QtyInBatch)

	for i := range s.config.QtyInBatch {
		logs[i] = s.generateLine()
	}

	return logs
}

func (s LogGenerationService) generateLine() string {
	ipAddress := generateIP()
	method := generateMethod()
	path := generatePath()
	status := generateStatus(s.config.ClientErrFreq, s.config.ServerErrFreq)
	latency := generateLatency()

	return fmt.Sprintf(
		"%v %v %v %v %v %vms",
		time.Now().Format(time.RFC3339),
		ipAddress,
		method,
		path,
		status,
		latency)
}

func generateIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Int63n(256), rand.Int63n(256), rand.Int63n(256), rand.Int63n(256))
}

func generateMethod() string {
	methods := []string{http.MethodPost, http.MethodGet, http.MethodPatch, http.MethodPut, http.MethodDelete}
	return methods[rand.Intn(len(methods))]
}

func generatePath() string {
	path := []string{
		"/",
		"/home",
		"/api/user",
		"/api/payment",
		"/api/logs",
		"/products",
		"/auth/login",
		"/health",
	}
	return path[rand.Intn(len(path))]
}

func generateStatus(clientErrFreq float64, serverErrFreq float64) int {
	r := rand.Float64()

	success := 1 - clientErrFreq - serverErrFreq

	if r <= success {
		return http.StatusOK
	}

	// Total errors freq
	totalErrors := clientErrFreq + serverErrFreq
	// Limiter for errors
	errorLimit := (r - success) / totalErrors

	switch {
	case errorLimit < (clientErrFreq / totalErrors):
		return http.StatusNotFound
	case errorLimit < ((serverErrFreq + clientErrFreq) / totalErrors):
		return http.StatusInternalServerError
	default:
		return 1
	}
}

func generateLatency() int64 {
	duration := time.Duration(rand.Intn(1000)) * time.Millisecond
	return duration.Milliseconds()
}
