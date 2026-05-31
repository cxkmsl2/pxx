package config

import "os"

type DBConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func (d DBConfig) DSN() string {
	return d.User + ":" + d.Pass + "@tcp(" + d.Host + ":" + d.Port + ")/" + d.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
}

type Config struct {
	AccountDB DBConfig
	ItemDB    DBConfig
	TradeDB   DBConfig
	FeedDB    DBConfig
	RedisAddr string
	KafkaBros string
	ServerPort string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		AccountDB: DBConfig{
			Host: env("ACCOUNT_DB_HOST", "127.0.0.1"),
			Port: env("ACCOUNT_DB_PORT", "3306"),
			User: env("DB_USER", "root"),
			Pass: env("DB_PASS", "root123456"),
			Name: env("ACCOUNT_DB_NAME", "pxx_account"),
		},
		ItemDB: DBConfig{
			Host: env("ITEM_DB_HOST", "127.0.0.1"),
			Port: env("ITEM_DB_PORT", "3306"),
			User: env("DB_USER", "root"),
			Pass: env("DB_PASS", "root123456"),
			Name: env("ITEM_DB_NAME", "pxx_item"),
		},
		TradeDB: DBConfig{
			Host: env("TRADE_DB_HOST", "127.0.0.1"),
			Port: env("TRADE_DB_PORT", "3306"),
			User: env("DB_USER", "root"),
			Pass: env("DB_PASS", "root123456"),
			Name: env("TRADE_DB_NAME", "pxx_trade"),
		},
		FeedDB: DBConfig{
			Host: env("FEED_DB_HOST", "127.0.0.1"),
			Port: env("FEED_DB_PORT", "3306"),
			User: env("DB_USER", "root"),
			Pass: env("DB_PASS", "root123456"),
			Name: env("FEED_DB_NAME", "pxx_feed"),
		},
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
