package tests

import (
	"ai-service/internal/infrastructure/s3"
	"context"
	"os"
	"testing"
)

func TestS3(t *testing.T) {
	S3AccessKey := os.Getenv("S3_ACCESS_KEY")
	S3SecretKey := os.Getenv("S3_SECRET_KEY")
	S3BucketName := os.Getenv("S3_BUCKET_NAME")
	endpoint := os.Getenv("S3_ENDPOINT")

	s3, err := s3.NewClient(S3AccessKey, S3SecretKey, S3BucketName, endpoint)
	if err != nil {
		t.Fatalf("Не удалось создать S3 клиента: %v", err)
	}

	_, err = s3.DownloadPDF(
		context.Background(),
		"https://storage.yandexcloud.net/trade-compass-reports/reports/x5/2025/6/x5_2025_6.pdf",
	)
	if err != nil {
		t.Fatalf("Не удалось скачать файл: %v", err)
	}
}
