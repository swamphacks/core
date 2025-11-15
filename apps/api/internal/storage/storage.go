package storage

import (
	"context"
)

type Storage interface {
	Store(ctx context.Context, bucketName, key string, data []byte, contentType *string) error
	Retrieve(ctx context.Context, bucketName, key string) ([]byte, error)
	Delete(ctx context.Context, bucketName, key string) error
	Close() error
}
