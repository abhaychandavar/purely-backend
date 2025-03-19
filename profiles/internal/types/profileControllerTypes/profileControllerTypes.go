package profileControllerTypes

type CreateProfileType struct {
	Lat *float64 `json:"lat"`
	Lng *float64 `json:"lng"`
}

type MediaElementType struct {
	MediaID *string `json:"mediaID"`
	Order   *int    `json:"order"`
}

type DatingPromptType struct {
	PromptId *string `json:"id"`
	Answer   *string `json:"answer"`
}

type Location struct {
	Lat *float64 `json:"lat"`
	Lng *float64 `json:"lng"`
}
type UpsertDatingProfileType struct {
	Name                   *string             `json:"name"`
	Age                    *int                `json:"age"`
	Gender                 *string             `json:"gender"`
	HereFor                *string             `json:"hereFor"`
	LookingFor             *string             `json:"lookingFor"`
	Bio                    *string             `json:"bio"`
	Prompts                *[]DatingPromptType `json:"prompts"`
	Media                  *[]MediaElementType `json:"media"`
	Lat                    *float64            `json:"lat"`
	Lng                    *float64            `json:"lng"`
	LocationLabel          *string             `json:"locationLabel"`
	Location               *Location           `json:"location"`
	PreferredMatchDistance *int                `json:"preferredMatchDistance"`
}

type GenerateMediaUploadSignedUrlType struct {
	FileName *string `json:"filename"`
	MimeType *string `json:"mimeType"`
	FileSize *int64  `json:"fileSize"`
	Purpose  *string `json:"purpose"`
}

type GenerateMultipartMediaUploadSignedUrls struct {
	FileName   string `json:"fileName"`
	MimeType   string `json:"mimeType"`
	AuthId     string `json:"authId"`
	FileSize   int64  `json:"fileSize"`
	Purpose    string `json:"purpose"`
	PartsCount int    `json:"partsCount"`
}

type CompleteMultipartUpload struct {
	UploadID string         `json:"uploadID"`
	FilePath string         `json:"filePath"`
	Parts    map[int]string `json:"parts"`
}
