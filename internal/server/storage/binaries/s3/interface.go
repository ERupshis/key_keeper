package s3

import (
	"context"

	"github.com/erupshis/key_keeper/internal/server/storage/binaries/models"
)

type BaseBucketManager interface {
	MakeBucket(ctx context.Context, bucketName string, location string) error
	ListBuckets(ctx context.Context) ([]models.Bucket, error)
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	RemoveBucket(ctx context.Context, bucketName string) error
}

type BaseObjectManager interface {
	PutObject(ctx context.Context, objectData *models.Object) error
	GetObject(ctx context.Context, objectShortData *models.ObjectName) (*models.Object, error)
	StatObject(ctx context.Context, objectShortData *models.ObjectName) (*models.ObjectStat, error)
	RemoveObject(ctx context.Context, objectShortData *models.ObjectName) error
	RemoveObjectsInBucket(ctx context.Context, bucketName string) error
	ListObjects(ctx context.Context, bucketName string) <-chan models.ObjectStat
}
