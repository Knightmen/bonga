package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	client *s3.Client
	bucket string
}

func NewS3Service() (*S3Service, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Get bucket name from environment variable
	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		return nil, fmt.Errorf("AWS_S3_BUCKET environment variable is not set")
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	return &S3Service{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *S3Service) GetPresignedURL(key string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	// Create presigned URL with 1-minute expiration
	presignedURL, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(time.Minute))

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return presignedURL.URL, nil
} 