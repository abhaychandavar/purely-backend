package mediaControllerTypes

type GenerateMediaUploadSignedUrlType struct {
	FileName    *string `json:"filename"`
	ContentType *string `json:"contentType"`
	FileSize    *int    `json:"fileSize"`
	Purpose     *string `json:"purpose"`
}

type GenerateMultipartMediaUploadSignedUrls struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	AuthId      string `json:"authId"`
	FileSize    int    `json:"fileSize"`
	Purpose     string `json:"purpose"`
}

type CompleteMultipartUpload struct {
	UploadID string         `json:"uploadID"`
	URL      string         `json:"url"`
	Parts    map[int]string `json:"parts"`
}
