package minio

import (
	"context"
	"fmt"

	"github.com/erupshis/key_keeper/internal/server/storage/binaries/models"
	"github.com/erupshis/key_keeper/internal/server/storage/binaries/s3"
	"github.com/minio/minio-go/v7"
)

var (
	_ s3.BaseBucketManager = (*BucketManager)(nil)
)

type BucketManager struct {
	client *minio.Client
}

func NewBucketManager(client *minio.Client) *BucketManager {
	return &BucketManager{client: client}
}

func (bm *BucketManager) MakeBucket(ctx context.Context, bucketName string, location string) error {
	err := bm.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := bm.client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			return fmt.Errorf("already added bucket: '%s'", bucketName)
		} else {
			return fmt.Errorf("add new bucket: %w", err)
		}
	}

	return nil
}

func (bm *BucketManager) ListBuckets(ctx context.Context) ([]models.Bucket, error) {
	buckets, err := bm.client.ListBuckets(ctx)
	if err != nil {
		return nil, fmt.Errorf("list buckets: %w", err)
	}

	var res []models.Bucket
	for _, bucket := range buckets {
		res = append(res, models.Bucket{
			Name: bucket.Name,
		})
	}

	return res, nil
}

func (bm *BucketManager) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	found, err := bm.client.BucketExists(ctx, bucketName)
	if err != nil {
		return false, fmt.Errorf("check bucket presence: %w", err)
	}

	return found, nil
}

func (bm *BucketManager) RemoveBucket(ctx context.Context, bucketName string) error {
	err := bm.client.RemoveBucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("remove bucket: %w", err)
	}

	return nil
}
