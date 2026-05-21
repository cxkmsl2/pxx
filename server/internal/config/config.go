package config

import "os"

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	RedisAddr  string
	KafkaBros  string
	ServerPort string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		DBHost:     env("DB_HOST", "127.0.0.1"),
		DBPort:     env("DB_PORT", "3306"),
		DBUser:     env("DB_USER", "pxx"),
		DBPass:     env("DB_PASS", "pxx123456"),
		DBName:     env("DB_NAME", "pxx"),
		RedisAddr:  env("REDIS_ADDR", "127.0.0.1:6379"),
		KafkaBros:  env("KAFKA_BROKERS", "127.0.0.1:9092"),
		ServerPort: env("SERVER_PORT", "8080"),
		JWTSecret:  env("JWT_SECRET", "pxx-secret-key-2025"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
