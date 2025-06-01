package config

type WorkerConfig struct {
	KafkaBrokers []string
	KafkaTopic   string
	DBDSN        string
}

func Load() *WorkerConfig {
	// TODO: Move to env
	return &WorkerConfig{
		KafkaBrokers: []string{"broker:29092"},
		KafkaTopic:   "raw-logs",
		DBDSN:        "app_user:app_pass@tcp(mysql:3306)/log_db?charset=utf8mb4&parseTime=True",
	}
}
