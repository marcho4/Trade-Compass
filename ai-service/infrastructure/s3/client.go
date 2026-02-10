package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3Client   *s3.Client
	bucketName string
}

func NewClient(accessKey, secretKey, bucketName, endpoint string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("ru-central1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &Client{
		s3Client:   s3Client,
		bucketName: bucketName,
	}, nil
}

func (c *Client) DownloadPDF(ctx context.Context, s3Path string) ([]byte, error) {
	key, err := extractKeyFromURL(s3Path, c.bucketName)
	if err != nil {
		return nil, err
	}

	log.Printf("Downloading PDF from S3: bucket=%s key=%s", c.bucketName, key)

	output, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}
	defer output.Body.Close()

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read S3 object body: %w", err)
	}

	log.Printf("Downloaded %d bytes from S3", len(data))
	return data, nil
}

func extractKeyFromURL(s3URL, bucketName string) (string, error) {
	parsed, err := url.Parse(s3URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse S3 URL: %w", err)
	}

	path := strings.TrimPrefix(parsed.Path, "/")
	path = strings.TrimPrefix(path, bucketName+"/")

	if path == "" {
		return "", fmt.Errorf("empty key extracted from URL: %s", s3URL)
	}

	return path, nil
}
