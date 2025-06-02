package config

import (
	"os"
	"strconv"
)

type Config struct {
	RPS           int
	QtyInBatch    int
	ClientErrFreq float64
	ServerErrFreq float64
	TargetApiUrl  string
}

func InitConfig() Config {
	rpsEnv := os.Getenv("RPS")
	rpsInt, err := strconv.Atoi(rpsEnv)
	if err != nil {
		rpsInt = 1
	}

	qtyEnv := os.Getenv("QTY_IN_BATCH")
	qtyInt, err := strconv.Atoi(qtyEnv)
	if err != nil {
		qtyInt = 100
	}

	cefEnv := os.Getenv("CLIENT_ERROR_FREQ")
	cefFloat, err := strconv.ParseFloat(cefEnv, 64)
	if err != nil {
		cefFloat = 0.2
	}

	sefEnv := os.Getenv("SERVER_ERROR_FREQ")
	sefFloat, err := strconv.ParseFloat(sefEnv, 64)
	if err != nil {
		sefFloat = 0.2
	}

	targetApiUrlEnv := os.Getenv("TARGET_API_URL")
	if targetApiUrlEnv == "" {
		targetApiUrlEnv = "http://localhost:8080/upload"
	}

	return Config{
		RPS:           rpsInt,
		QtyInBatch:    qtyInt,
		ClientErrFreq: cefFloat,
		ServerErrFreq: sefFloat,
		TargetApiUrl:  targetApiUrlEnv,
	}
}
