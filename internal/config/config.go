package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort              string
	DBURL                string
	KafkaBrokers         string
	KafkaTopicRegCreated string
	KafkaTopicPayUploaded string
	KafkaTopicPayVerified string
	KafkaTopicRegConfirmed string
	KafkaTopicRegCancelled string
	KafkaTopicEventStatus string
}

func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:              getEnv("APP_PORT", "3003"),
		DBURL:                getEnv("DB_URL", "postgres://regpay:regpay@localhost:5435/regpay_db?sslmode=disable"),
		KafkaBrokers:         getEnv("KAFKA_BROKERS", "localhost:19093"),
		KafkaTopicRegCreated: getEnv("KAFKA_TOPIC_REG_CREATED", "registration.created"),
		KafkaTopicPayUploaded: getEnv("KAFKA_TOPIC_PAY_UPLOADED", "payment.uploaded"),
		KafkaTopicPayVerified: getEnv("KAFKA_TOPIC_PAY_VERIFIED", "payment.verified"),
		KafkaTopicRegConfirmed: getEnv("KAFKA_TOPIC_REG_CONFIRMED", "registration.confirmed"),
		KafkaTopicRegCancelled: getEnv("KAFKA_TOPIC_REG_CANCELLED", "registration.cancelled"),
		KafkaTopicEventStatus: getEnv("KAFKA_TOPIC_EVENT_STATUS", "event.status.changed"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DBURL == "" {
		return fmt.Errorf("DB_URL is required")
	}
	if c.KafkaBrokers == "" {
		return fmt.Errorf("KAFKA_BROKERS is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

