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

type CompletedMultipartUploadResponseType struct {
	URL      string
	Path     string
	Domain   string
	FileSize int64
}

type UploadResType struct {
	URL    string
	Path   string
	Domain string
}
