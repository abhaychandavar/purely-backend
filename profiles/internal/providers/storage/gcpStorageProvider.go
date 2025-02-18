package storage

import (
	"context"
	"profiles/internal/config"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	clientInstance *storage.Client
	once           sync.Once
)

type GCPStorageProvider struct{}

func (provider *GCPStorageProvider) getClient(ctx *context.Context) (*storage.Client, error) {
	var err error
	once.Do(func() {
		var opts option.ClientOption
		if config.GetConfig().Env == "development" {
			opts = option.WithCredentialsFile(config.GetConfig().GoogleServiceJsonFilePath)
		}
		clientInstance, err = storage.NewClient(*ctx, opts)
	})

	return clientInstance, err
}

type InitiateMultipartUploadResult struct {
	UploadID string `xml:"UploadId"`
}

func (provider *GCPStorageProvider) GenerateSignedUrl(ctx *context.Context, bucket string, filePath string, MimeType string, fileSize int64) (*UploadSignedUrl, error) {
	client, err := provider.getClient(ctx)
	if err != nil {
		return nil, err
	}
	storageBucket := client.Bucket(bucket)
	expiry := time.Now().Add(10 * time.Minute)
	signedUrl, err := storageBucket.SignedURL(filePath, &storage.SignedURLOptions{
		Expires: expiry,
	})
	if err != nil {
		return nil, err
	}
	return &UploadSignedUrl{
		Bucket:    bucket,
		FilePath:  filePath,
		SignedUrl: signedUrl,
		Expires:   expiry,
	}, nil
}
