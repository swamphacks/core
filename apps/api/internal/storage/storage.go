package storage

import (
	"context"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

type Storage interface {
	Store(ctx context.Context, bucketName, key string, data []byte, contentType *string) error
	Retrieve(ctx context.Context, bucketName, key string) ([]byte, error)
	Delete(ctx context.Context, bucketName, key string) error
	PresignGetObject(ctx context.Context, bucketName, key string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error)
	Close() error
}
