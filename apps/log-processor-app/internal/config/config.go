package config

type WorkerConfig struct {
	KafkaBrokers []string
	KafkaTopic   string
	DBDSN        string
	RedisAddr    string
}

func Load() *WorkerConfig {
	return &WorkerConfig{
		KafkaBrokers: []string{"localhost:9092"},
		KafkaTopic:   "raw-logs",
		DBDSN:        "app_user:app_pass@tcp(localhost:3306)/log_db?charset=utf8mb4&parseTime=True",
		RedisAddr:    "reeee",
	}
}
