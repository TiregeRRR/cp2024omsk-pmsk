package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/gulldan/cp2024omsk-pmsk/config"
	"github.com/minio/minio-go/v7/pkg/credentials"

	minio "github.com/minio/minio-go/v7"
)

type MinioClient struct {
	client *minio.Client
}

func NewMinioClient(opts *config.Config) (*MinioClient, error) {
	minioClient, err := minio.New(
		opts.MinioEndpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(opts.MinioAccessKey, opts.MinioSecretAccessKey, ""),
			Secure: false,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to create new minio client: %w", err)
	}

	return &MinioClient{
		client: minioClient,
	}, nil
}

func (s *MinioClient) UploadFile(ctx context.Context, data io.Reader, dataSize int64, objectName, bucketName string) error {
	if exist := s.isBucketExist(ctx, bucketName); !exist {
		if err := s.makeBucket(ctx, bucketName); err != nil {
			return fmt.Errorf("failed to make bucket when upload file to s3: %w", err)
		}
	}

	_, err := s.client.PutObject(ctx, bucketName, objectName, data, dataSize, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to put object in s3: %w", err)
	}

	return nil
}

func (s *MinioClient) isBucketExist(ctx context.Context, bucketName string) bool {
	exists, errBucketExists := s.client.BucketExists(ctx, bucketName)
	if errBucketExists == nil && exists {
		return true
	}

	return false
}

func (s *MinioClient) makeBucket(ctx context.Context, bucketName string) error {
	if err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("failed to make bucket:%w", err)
	}

	return nil
}

func (s *MinioClient) GetAudioBucket() string {
	return "audio"
}

func (s *MinioClient) DownloadFile(ctx context.Context, objectName, bucketName string) (io.Reader, error) {
	reader, err := s.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from s3: %w", err)
	}

	return reader, nil
}
