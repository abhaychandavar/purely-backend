package storage

import "time"

type GenerateSignedURLsForPartsResType struct {
	SignedUrls map[int]string
	Expiry     time.Time
}

type StorageProvider interface {
	GenerateSignedUrl(bucket string, filePath string, MimeType string, fileSize int64) (*UploadSignedUrl, error)
	InitiateMultipartUpload(bucket string, filePath string, MimeType string, fileSize int64) (*InitiateMultipartUpload, error)
	GenerateSignedURLsForParts(bucket string, filePath string, uploadID string, partsCount int) (*GenerateSignedURLsForPartsResType, error)
	CompleteMultipartUpload(bucket string, uploadID string, filePath string, parts map[int]string) (*string, error)
}
