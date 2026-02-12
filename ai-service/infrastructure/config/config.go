package config

import "os"

type Config struct {
	APIKey              string
	GeminiAPIKey        string
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
}

func Load() *Config {
	return &Config{
		APIKey:              getEnv("AI_SERVICE_API_KEY", ""),
		GeminiAPIKey:        getEnv("GEMINI_API_KEY", ""),
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
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
