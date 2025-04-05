package mediaServiceTypes

type GenerateMediaUploadSignedUrlType struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	AuthId      string `json:"authId"`
	FileSize    int    `json:"fileSize"`
	Purpose     string `json:"purpose"`
}

type GenerateMediaUploadSignedUrlResType struct {
	SignedUrl string `json:"signedUrl"`
	Expiry    int64  `json:"expiry"`
}

type GenerateMultipartUploadUrlsType struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	AuthId      string `json:"authId"`
	FileSize    int    `json:"fileSize"`
	Purpose     string `json:"purpose"`
	PartsCount  int    `json:"partsCount"`
}

type GenerateMultipartUploadUrlsResType struct {
	SignedUrls map[int]string `json:"signedUrls"`
	Expiry     int64          `json:"expiry"`
	UploadID   string         `json:"uploadID"`
	FilePath   string         `json:"filePath"`
	PartsCount int            `json:"partsCount"`
	URL        string         `json:"url"`
}

type CompleteMultipartUploadType struct {
	UploadID string         `json:"uploadID"`
	URL      string         `json:"url"`
	Parts    map[int]string `json:"parts"`
}

type CompleteMultipartUploadResType struct {
	URL string `json:"url"`
	ID  string `json:"id"`
}
