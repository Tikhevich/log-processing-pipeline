package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/muonsoft/validation/validate"
	"worker.go/internal/models"
)

func ParseLog(line string) (*models.LogEntry, error) {
	line = strings.Trim(line, "\"")
	splitedStr := strings.Fields(line)
	if len(splitedStr) != 6 {
		return nil, fmt.Errorf("Invalid log format")
	}

	timeStamp, err := time.Parse(time.RFC3339, splitedStr[0])
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	isIp := checkIsIp(splitedStr[1])
	if isIp == false {
		return nil, fmt.Errorf("invalid ip address")
	}

	isRegisteredHTTPMethod := checkIsRegisteredHTTPMethod(splitedStr[2])
	if isRegisteredHTTPMethod == false {
		return nil, fmt.Errorf("invalid http method")
	}

	isPath := strings.HasPrefix(splitedStr[3], "/")
	if isPath == false {
		return nil, fmt.Errorf("invalid path format")
	}

	status, err := strconv.Atoi(splitedStr[4])
	if err != nil {
		return nil, fmt.Errorf("invalid status code: %w", err)
	}

	latencyStr := strings.TrimSuffix(splitedStr[5], "ms")
	latency, err := strconv.Atoi(latencyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid latency: %w", err)
	}

	return &models.LogEntry{
		Timestamp: timeStamp,
		IP:        splitedStr[1],
		Method:    splitedStr[2],
		Path:      splitedStr[3],
		Status:    status,
		LatencyMs: latency,
	}, nil
}

func ParseLogs(data []byte) ([]models.LogEntry, error) {
	var lines []string
	err := json.Unmarshal(data, &lines)
	if err != nil {
		return nil, err
	}

	entries := make([]models.LogEntry, 0, len(lines))

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		entry, err := ParseLog(line)
		if err != nil {
			fmt.Println(err)
			continue
		}
		entries = append(entries, *entry)
	}

	return entries, nil
}

func checkIsRegisteredHTTPMethod(method string) bool {
	registeredMethods := []string{http.MethodPost, http.MethodGet, http.MethodPatch, http.MethodPut, http.MethodDelete}
	return slices.Contains(registeredMethods, method)
}

func checkIsIp(stringIp string) bool {
	err := validate.IP(stringIp)
	return err == nil
}
