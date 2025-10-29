package storage

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"
)

type R2Client struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	logger        zerolog.Logger
}

// NewR2Client initializes a new R2Client with the provided bucket name and logger.
func NewR2Client(accountId, accessKey, secretkey string, logger zerolog.Logger) (*R2Client, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretkey, "")),
		awsConfig.WithRegion("auto"),
	)

	if err != nil {
		logger.Err(err).Msg("Failed to load AWS configuration")
		return nil, err
	}

	// Config stuff
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	return &R2Client{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		logger:        logger.With().Str("component", "s3_client").Logger(),
	}, nil
}

func (c *R2Client) Store(ctx context.Context, bucketName, key string, data []byte, contentType *string) error {
	r := bytes.NewReader(data)

	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        r,
		ContentType: contentType,
	})
	if err != nil {
		c.logger.Err(err).Msgf("Failed to upload object to S3 with key %s", key)
		return err
	}

	return nil
}

func (c *R2Client) Retrieve(ctx context.Context, bucketName, key string) ([]byte, error) {
	result, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		c.logger.Err(err).Msgf("Failed to retrieve object from S3 with key %s", key)
		return nil, err
	}

	defer result.Body.Close()
	data := new(bytes.Buffer)
	if _, err := data.ReadFrom(result.Body); err != nil {
		c.logger.Err(err).Msgf("Failed to read object body from S3 with key %s", key)
		return nil, err
	}

	return data.Bytes(), nil
}

func (c *R2Client) Delete(ctx context.Context, bucketName, key string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		c.logger.Err(err).Msgf("Failed to delete object from S3 with key %s", key)
		return err
	}

	return nil
}

// https://github.com/awsdocs/aws-doc-sdk-examples/blob/309de24d867a2b3f01d4da8018ad8173243556af/gov2/s3/actions/presigner.go#L33
func (c *R2Client) PresignGetObject(ctx context.Context, bucketName, key string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := c.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})

	if err != nil {
		c.logger.Err(err).Msgf("Failed to create presigned get object to S3 with key %s", key)
		return nil, err
	}

	return request, err
}

func (c *R2Client) Close() error {
	// R2Client does not require explicit closure, but you can implement any cleanup logic if needed.
	// This is because the underlying S3 client does not maintain persistent connections.
	return nil
}
