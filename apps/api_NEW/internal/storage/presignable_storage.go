package storage

import (
	"context"
	"net/http"
)

type PresignedRequest struct {
	URL     string
	Method  string
	Headers http.Header
}

type PresignableStorage interface {
	Storage
	PresignGetObject(ctx context.Context, bucketName, key string, lifetimeSecs int64) (*PresignedRequest, error)
}
