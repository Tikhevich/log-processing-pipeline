package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SendLogService struct {
	logGenerator *LogGenerationService
}

func NewSendLogService(logGenerator *LogGenerationService) *SendLogService {
	return &SendLogService{logGenerator: logGenerator}
}

func (s SendLogService) SendLogs(url string) {
	logs := s.logGenerator.GenerateLogs()

	jsonValue, err := json.Marshal(logs)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("HTTP error:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Println(string(body))
}
