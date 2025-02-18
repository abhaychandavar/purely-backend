package storage

import "context"

type StorageProvider interface {
	GenerateSignedUrl(ctx *context.Context, bucket string, filePath string, MimeType string, fileSize int64) (*UploadSignedUrl, error)
}
