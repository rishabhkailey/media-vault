package utils

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func CreateBucketIfMissing(ctx context.Context, client minio.Client, bucketName string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}
