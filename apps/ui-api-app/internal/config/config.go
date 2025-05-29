package config

type Config struct {
	DBDSN     string
	RedisAddr string
}

func Load() *Config {
	// TODO: Update to env
	return &Config{
		DBDSN:     "app_user:app_pass@tcp(mysql:3306)/log_db?charset=utf8mb4&parseTime=True",
		RedisAddr: "redis:6379",
	}
}
