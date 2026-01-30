package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	WorkerCount int
	RequestTimeout int
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		WorkerCount:    getEnvInt("WORKER_COUNT", 50),
		RequestTimeout: getEnvInt("REQUEST_TIMEOUT", 5),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	
	return defaultVal
}
