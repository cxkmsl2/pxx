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
	// ETCD 服务发现
	ETCDEndpoints string

	// 各微服务数据库
	TradeDB   DBConfig
	AccountDB DBConfig
	ItemDB    DBConfig
	FeedDB    DBConfig

	// Redis
	RedisAddr string

	// Kafka (仅 Feed 服务使用)
	KafkaBros string

	// 本服务端口
	ServerPort string

	// JWT
	JWTSecret string
}

func Load() *Config {
	return &Config{
		ETCDEndpoints: env("ETCD_ENDPOINTS", "127.0.0.1:2379"),

		TradeDB: DBConfig{
			Host: env("TRADE_DB_HOST", "127.0.0.1"),
			Port: env("TRADE_DB_PORT", "3306"),
			User: env("TRADE_DB_USER", "pxx"),
			Pass: env("TRADE_DB_PASS", "pxx123456"),
			Name: env("TRADE_DB_NAME", "pxx_trade"),
		},
		AccountDB: DBConfig{
			Host: env("ACCOUNT_DB_HOST", "127.0.0.1"),
			Port: env("ACCOUNT_DB_PORT", "3306"),
			User: env("ACCOUNT_DB_USER", "pxx"),
			Pass: env("ACCOUNT_DB_PASS", "pxx123456"),
			Name: env("ACCOUNT_DB_NAME", "pxx_account"),
		},
		ItemDB: DBConfig{
			Host: env("ITEM_DB_HOST", "127.0.0.1"),
			Port: env("ITEM_DB_PORT", "3306"),
			User: env("ITEM_DB_USER", "pxx"),
			Pass: env("ITEM_DB_PASS", "pxx123456"),
			Name: env("ITEM_DB_NAME", "pxx_item"),
		},
		FeedDB: DBConfig{
			Host: env("FEED_DB_HOST", "127.0.0.1"),
			Port: env("FEED_DB_PORT", "3306"),
			User: env("FEED_DB_USER", "pxx"),
			Pass: env("FEED_DB_PASS", "pxx123456"),
			Name: env("FEED_DB_NAME", "pxx_feed"),
		},

		RedisAddr:  env("REDIS_ADDR", "127.0.0.1:6379"),
		KafkaBros:  env("KAFKA_BROKERS", "127.0.0.1:9092"),
		ServerPort: env("SERVER_PORT", "8080"),
		JWTSecret:  env("JWT_SECRET", "pxx-jwt-secret-2026"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
