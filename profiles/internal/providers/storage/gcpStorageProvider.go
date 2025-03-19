package storage

import (
	"context"
	"fmt"
	"profiles/internal/config"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type gcpStorageProvider struct {
	clientInstance *storage.Client
	once           sync.Once
}

func (provider *gcpStorageProvider) NewGCPStorageProvider(ctx context.Context) (*gcpStorageProvider, error) {
	var err error
	provider.once.Do(func() {
		var opts option.ClientOption
		if config.GetConfig().Env == "development" {
			opts = option.WithCredentialsFile(config.GetConfig().FirebaseConfigPath)
		}
		provider.clientInstance, err = storage.NewClient(ctx, opts)
		if err != nil {
			panic(err)
		}
	})
	return provider, err
}

type InitiateMultipartUploadResult struct {
	UploadID string `xml:"UploadId"`
}

func (provider *gcpStorageProvider) GenerateSignedUrl(ctx context.Context, bucket string, filePath string, MimeType string, fileSize int64) (*UploadSignedUrl, error) {
	client := provider.clientInstance
	storageBucket := client.Bucket(bucket)
	expiry := time.Now().Add(10 * time.Minute)
	signedUrl, err := storageBucket.SignedURL(filePath, &storage.SignedURLOptions{
		Expires: expiry,
		Method:  "PUT",
		Headers: []string{
			fmt.Sprintf("Content-Type: %s", MimeType),
			"Access-Control-Allow-Origin: *",
			"Access-Control-Allow-Methods: PUT",
		},
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
