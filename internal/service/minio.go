package service

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jonathanCaamano/portfolio-gateway/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	client *minio.Client
	bucket string
}

func NewMinioService(cfg *config.Config) *MinioService {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		panic(fmt.Sprintf("minio init: %v", err))
	}
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.MinioBucket)
	if err != nil {
		panic(fmt.Sprintf("minio bucket check: %v", err))
	}
	if !exists {
		if err = client.MakeBucket(ctx, cfg.MinioBucket, minio.MakeBucketOptions{}); err != nil {
			panic(fmt.Sprintf("minio make bucket: %v", err))
		}
	}
	return &MinioService{client, cfg.MinioBucket}
}

func (s *MinioService) Upload(objectKey string, data []byte, contentType string) (string, error) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	_, err := s.client.PutObject(context.Background(), s.bucket, objectKey,
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	return s.PresignedURL(objectKey)
}

func (s *MinioService) PresignedURL(objectKey string) (string, error) {
	u, err := s.client.PresignedGetObject(
		context.Background(), s.bucket, objectKey,
		7*24*time.Hour, url.Values{})
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
