package config

import (
	"fmt"
	"os"
)

type Config struct {
	APIKey              string
	GeminiAPIKey        string
	GeminiProxyURL      string
	S3AccessKey         string
	S3SecretKey         string
	S3BucketName        string
	S3Endpoint          string
	ParserURL           string
	FinancialDataURL    string
	FinancialDataAPIKey string
	Port                string
	KafkaURL            string
	KafkaTopic          string
	PostgresURL         string
}

func Load() *Config {
	return &Config{
		APIKey:              getEnv("AI_SERVICE_API_KEY", ""),
		GeminiAPIKey:        getEnv("GEMINI_API_KEY", ""),
		GeminiProxyURL:      getEnv("GEMINI_PROXY_URL", ""),
		S3AccessKey:         getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey:         getEnv("S3_SECRET_KEY", ""),
		S3BucketName:        getEnv("S3_BUCKET_NAME", ""),
		S3Endpoint:          getEnv("S3_ENDPOINT", "https://storage.yandexcloud.net"),
		ParserURL:           getEnv("PARSER_URL", "http://parser:8081"),
		FinancialDataURL:    getEnv("FINANCIAL_DATA_URL", "http://financial-data:8082"),
		FinancialDataAPIKey: getEnv("FINANCIAL_DATA_API_KEY", ""),
		Port:                getEnv("PORT", "8083"),
		KafkaURL:            getEnv("KAFKA_URL", "kafka:9092"),
		KafkaTopic:          getEnv("KAFKA_TOPIC", "ai-analyze-tasks"),
		PostgresURL:         getEnv("POSTGRES_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API_KEY is not set")
	}

	if c.GeminiAPIKey == "" {
		return fmt.Errorf("GeminiAPIKey is not set")
	}

	if c.GeminiProxyURL == "" {
		return fmt.Errorf("GeminiProxyURL is not set")
	}

	if c.PostgresURL == "" {
		return fmt.Errorf("POSTGRES_URL is not set")
	}

	if c.FinancialDataAPIKey == "" {
		return fmt.Errorf("FINANCIAL_DATA_API_KEY is not set")
	}

	if c.S3AccessKey == "" {
		return fmt.Errorf("S3AccessKey is not set")
	}

	if c.S3BucketName == "" {
		return fmt.Errorf("S3BucketName is not set")
	}

	if c.S3SecretKey == "" {
		return fmt.Errorf("S3SecretKey is not set")
	}
	return nil
}
