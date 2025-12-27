package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	DBUrl       string
	WorkerCount int
	BatchSize   int
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBUrl:       getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/sijil"),
		WorkerCount: getEnvAsInt("WORKER_COUNT", 25),
		BatchSize:   getEnvAsInt("BATCH_SIZE", 10000),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}
