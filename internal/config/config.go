package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddress     string
	MySQLDSN        string
	RedisAddress    string
	RedisPassword   string
	RedisDB         int
	KafkaBrokers    []string
	KafkaTopic      string
	KafkaConsumerID string
}

func Load() (Config, error) {
	address := os.Getenv("OPENRESTY_HTTP_ADDR")
	if address == "" {
		address = ":8081"
	}
	jdbcURL := strings.TrimPrefix(os.Getenv("OPENRESTY_DB_URL"), "jdbc:")
	if jdbcURL == "" {
		return Config{}, fmt.Errorf("OPENRESTY_DB_URL is required")
	}
	parsed, err := url.Parse(jdbcURL)
	if err != nil {
		return Config{}, fmt.Errorf("parse OPENRESTY_DB_URL: %w", err)
	}
	user := url.QueryEscape(os.Getenv("OPENRESTY_DB_USERNAME"))
	password := url.QueryEscape(os.Getenv("OPENRESTY_DB_PASSWORD"))
	if user == "" {
		return Config{}, fmt.Errorf("OPENRESTY_DB_USERNAME is required")
	}
	if strings.TrimSpace(os.Getenv("OPENRESTY_ADMIN_PASSWORD")) == "" {
		return Config{}, fmt.Errorf("OPENRESTY_ADMIN_PASSWORD is required")
	}
	if len(strings.TrimSpace(os.Getenv("OPENRESTY_DATA_KEY"))) != 64 {
		return Config{}, fmt.Errorf("OPENRESTY_DATA_KEY must be 32 random bytes encoded as 64 hex characters")
	}
	brokers := []string{}
	for _, broker := range strings.Split(os.Getenv("OPENRESTY_KAFKA_BOOTSTRAP_SERVERS"), ",") {
		if broker = strings.TrimSpace(broker); broker != "" {
			brokers = append(brokers, broker)
		}
	}
	topic := strings.TrimSpace(os.Getenv("OPENRESTY_KAFKA_TOPIC"))
	if topic == "" {
		topic = "dtl-602-openresty-plus"
	}
	consumerID := strings.TrimSpace(os.Getenv("OPENRESTY_KAFKA_CONSUMER_GROUP"))
	if consumerID == "" {
		consumerID = "openresty-plus-control-plane-v1"
	}
	redisAddress := strings.TrimSpace(os.Getenv("OPENRESTY_REDIS_ADDR"))
	if redisAddress == "" {
		host := strings.TrimSpace(os.Getenv("OPENRESTY_REDIS_HOST"))
		port := strings.TrimSpace(os.Getenv("OPENRESTY_REDIS_PORT"))
		if host == "" {
			host = "127.0.0.1"
		}
		if port == "" {
			port = "6379"
		}
		redisAddress = net.JoinHostPort(host, port)
	}
	redisDB := 0
	if raw := strings.TrimSpace(os.Getenv("OPENRESTY_REDIS_DB")); raw != "" {
		parsedDB, parseErr := strconv.Atoi(raw)
		if parseErr != nil || parsedDB < 0 {
			return Config{}, fmt.Errorf("OPENRESTY_REDIS_DB must be a non-negative integer")
		}
		redisDB = parsedDB
	}
	return Config{HTTPAddress: address, MySQLDSN: fmt.Sprintf("%s:%s@tcp(%s)%s?parseTime=true&charset=utf8mb4&multiStatements=true", user, password, parsed.Host, parsed.EscapedPath()), RedisAddress: redisAddress, RedisPassword: os.Getenv("OPENRESTY_REDIS_PASSWORD"), RedisDB: redisDB, KafkaBrokers: brokers, KafkaTopic: topic, KafkaConsumerID: consumerID}, nil
}
