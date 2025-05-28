package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type Config struct {
	RPS           int
	ClientErrFreq float64
	ServerErrFreq float64
	TargetApiUrl  string
}

func main() {

	config := initFlags()

	ticker := time.NewTicker(time.Second / time.Duration(config.RPS))
	defer ticker.Stop()

	for range ticker.C {

		var logs = make([]string, 100)
		for i := range 100 {
			logs[i] = generateLine(config)
		}

		sendLogs(logs)
	}
}

func sendLogs(logs []string) {
	jsonValue, err := json.Marshal(logs)
	response, err := http.Post("http://localhost:8080/upload", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	response.Body.Close()

	fmt.Println(string(body))
}

func generateLine(config Config) string {
	ipAddress := generateIP()
	method := generateMethod()
	path := generatePath()
	status := generateStatus(config.ClientErrFreq, config.ServerErrFreq)
	latency := generateLatency()

	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

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

func initFlags() Config {
	var config Config

	flag.IntVar(&config.RPS, "rps", 10, "Requests per second")
	flag.StringVar(&config.TargetApiUrl, "target-api", "http://localhost:8080/upload", "Target api for uploading logs")
	flag.Float64Var(&config.ClientErrFreq, "cef", 0.1, "help message for clientErrFreq flag")
	flag.Float64Var(&config.ServerErrFreq, "sef", 0.1, "help message for serverErrFreq flag")

	flag.Parse()

	return config
}
