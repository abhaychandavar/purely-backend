package storage

import "time"

type GenerateSignedURLsForPartsResType struct {
	SignedUrls map[int]string
	Expiry     time.Time
	PartsCount int
	URL        string
}

type StorageProvider interface {
	GenerateSignedUrl(bucket string, filePath string, fileName string, contentType string, fileSize int) (*UploadSignedUrl, error)
	InitiateMultipartUpload(bucket string, filePath string, fileName string, contentType string, fileSize int) (*InitiateMultipartUpload, error)
	GenerateSignedURLsForParts(bucket string, filePath string, fileName string, uploadID string, contentType string, fileSize int) (*GenerateSignedURLsForPartsResType, error)
	CompleteMultipartUpload(bucket string, uploadID string, filePath string, fileName string, contentType string, parts map[int]string) (*CompletedMultipartUploadResponseType, error)
	UploadFile(signedUrls map[int]string, data []byte, partSize int, contentType string) (map[int]string, error)
}
