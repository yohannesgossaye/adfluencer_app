package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithy "github.com/aws/smithy-go"
)

type Config struct {
	Endpoint     string
	Region       string
	AccessKeyID  string
	SecretKey    string
	Bucket       string
	PublicURL    string
	UsePathStyle bool
}

type S3Storage struct {
	client    *s3.Client
	uploader  *manager.Uploader
	bucket    string
	publicURL string
}

func NewS3(cfg Config) (*S3Storage, error) {
	if !cfg.Enabled() {
		return nil, errors.New("s3 is not configured")
	}

	region := cfg.Region
	if region == "" {
		region = "us-east-1"
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("load s3 config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(strings.TrimRight(cfg.Endpoint, "/"))
		o.UsePathStyle = cfg.UsePathStyle
	})

	return &S3Storage{
		client:    client,
		uploader:  manager.NewUploader(client),
		bucket:    cfg.Bucket,
		publicURL: strings.TrimRight(cfg.PublicURL, "/"),
	}, nil
}

func (c Config) Enabled() bool {
	return strings.TrimSpace(c.Endpoint) != "" &&
		strings.TrimSpace(c.AccessKeyID) != "" &&
		strings.TrimSpace(c.SecretKey) != "" &&
		strings.TrimSpace(c.Bucket) != ""
}

func (s *S3Storage) Enabled() bool {
	return s != nil && s.client != nil
}

func (s *S3Storage) PublicURLBase() string {
	return s.publicURL
}

func (s *S3Storage) UploadBytes(ctx context.Context, objectKey string, data io.Reader, size int64, contentType string) (string, error) {
	objectKey = strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	if objectKey == "" {
		return "", errors.New("object key is required")
	}

	_, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(objectKey),
		Body:          data,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("upload object %s: %w", objectKey, err)
	}

	return objectKey, nil
}

func (s *S3Storage) Check(ctx context.Context, objectKey string) (bool, error) {
	objectKey = strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	if objectKey == "" {
		return false, errors.New("object key is required")
	}

	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})
	if err == nil {
		return true, nil
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "NotFound", "NoSuchKey", "404":
			return false, nil
		}
	}

	return false, fmt.Errorf("check object %s: %w", objectKey, err)
}

func (s *S3Storage) Delete(ctx context.Context, objectKey string) error {
	objectKey = strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	if objectKey == "" {
		return errors.New("object key is required")
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("delete object %s: %w", objectKey, err)
	}
	return nil
}

func (s *S3Storage) GetFileURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	objectKey = strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	if objectKey == "" {
		return "", errors.New("object key is required")
	}

	if expiry > 0 {
		_ = ctx
		// Presigned URLs can be added later for private buckets.
	}

	escaped := strings.ReplaceAll(objectKey, " ", "%20")
	fileURL, err := url.JoinPath(s.publicURL, escaped)
	if err != nil {
		return "", fmt.Errorf("build file url: %w", err)
	}
	return fileURL, nil
}
