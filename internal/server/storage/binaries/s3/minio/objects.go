package minio

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/erupshis/key_keeper/internal/common/utils/deferutils"
	"github.com/erupshis/key_keeper/internal/server/storage/binaries/models"
	"github.com/erupshis/key_keeper/internal/server/storage/binaries/s3"
	"github.com/minio/minio-go/v7"
)

var (
	_ s3.BaseObjectManager = (*ObjectManager)(nil)
)

type ObjectManager struct {
	*minio.Client
}

func NewObjectManager(client *minio.Client) *ObjectManager {
	return &ObjectManager{Client: client}
}

func (om *ObjectManager) PutObject(ctx context.Context, objectData *models.Object) error {
	_, err := om.Client.PutObject(ctx,
		objectData.Bucket,
		objectData.Name,
		strings.NewReader(objectData.Data),
		objectData.Size,
		minio.PutObjectOptions{ContentType: objectData.ContentType},
	)

	if err != nil {
		return fmt.Errorf("put object: '%w'", err)
	}

	return nil
}

func (om *ObjectManager) GetObject(ctx context.Context, objectShortData *models.ObjectName) (*models.Object, error) {
	object, err := om.Client.GetObject(ctx,
		objectShortData.Bucket,
		objectShortData.Name,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	defer deferutils.ExecSilent(object.Close)

	res := models.Object{Name: objectShortData.Name}

	stat, err := object.Stat()
	if err != nil {
		return nil, fmt.Errorf("get object stat: %w", err)
	}

	res.Size = stat.Size
	res.ContentType = stat.ContentType

	buf := bytes.Buffer{}

	if _, err = buf.ReadFrom(object); err != nil {
		return nil, fmt.Errorf("copy object's data to result: %w", err)
	}

	res.Data = buf.String()
	return &res, nil
}

func (om *ObjectManager) StatObject(ctx context.Context, objectShortData *models.ObjectName) (*models.ObjectStat, error) {
	stat, err := om.Client.StatObject(ctx,
		objectShortData.Bucket,
		objectShortData.Name,
		minio.StatObjectOptions{},
	)

	if err != nil {
		return nil, fmt.Errorf("get object's stat: %w", err)
	}

	res := models.ObjectStat{
		ETag:         stat.ETag,
		Key:          stat.Key,
		LastModified: stat.LastModified,
		Size:         stat.Size,
		ContentType:  stat.ContentType,
		Expires:      stat.Expires,
		Metadata:     stat.Metadata,
		VersionID:    stat.VersionID,
	}
	return &res, nil
}

func (om *ObjectManager) RemoveObject(ctx context.Context, objectShortData *models.ObjectName) error {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	err := om.Client.RemoveObject(ctx,
		objectShortData.Bucket,
		objectShortData.Name, opts,
	)

	if err != nil {
		return fmt.Errorf("remove object: %w", err)
	}

	return nil
}

func (om *ObjectManager) RemoveObjectsInBucket(ctx context.Context, bucketName string) error {
	objectsCh := make(chan minio.ObjectInfo)

	listOpts := minio.ListObjectsOptions{
		Recursive: true,
	}

	go func() {
		defer close(objectsCh)
		for object := range om.Client.ListObjects(ctx, bucketName, listOpts) {
			if object.Err != nil {
				// TODO: handle err.
				break
			}
			objectsCh <- object
		}
	}()

	removeOpts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	var err error
	for removeErr := range om.Client.RemoveObjects(ctx, bucketName, objectsCh, removeOpts) {
		err = fmt.Errorf("remove object '%s': %w", removeErr.ObjectName, removeErr.Err)
	}

	return err
}

func (om *ObjectManager) ListObjects(ctx context.Context, bucketName string) <-chan models.ObjectStat {
	objectCh := om.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	resCh := make(chan models.ObjectStat, 1)
	go func() {
		for object := range objectCh {
			if object.Err != nil {
				// TODO: need to handle somehow.
				break
			}

			tmpObj := <-objectCh

			resCh <- models.ObjectStat{
				Key:       tmpObj.Key,
				VersionID: tmpObj.VersionID,
			}
		}

		close(resCh)
	}()

	return resCh
}
