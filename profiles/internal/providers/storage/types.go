package storage

import "time"

type UploadSignedUrl struct {
	Bucket    string
	FilePath  string
	SignedUrl string
	Expires   time.Time
}

type InitiateMultipartUpload struct {
	Bucket   string
	FilePath string
	UploadId string
}
