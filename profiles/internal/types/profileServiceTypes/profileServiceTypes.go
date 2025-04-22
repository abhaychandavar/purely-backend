package profileServiceTypes

import (
	"profiles/internal/database/models"
	"time"
)

type CreateProfileType struct {
	AuthId   *string  `json:"authId"`
	Lat      *float64 `json:"lat"`
	Lng      *float64 `json:"lng"`
	Category *string  `json:"category"`
}

type GetProfileType struct {
	Category *string `json:"category"`
	AuthId   *string `json:"authId"`
}

type GetProfileLayoutType struct {
	Category *string `json:"category"`
}

type MediaElementType struct {
	MediaID        string `json:"mediaID"`
	Order          int    `json:"order"`
	BlurredImageID string `json:"blurredImageID"`
}

type DatingPromptType struct {
	PromptId string `json:"id"`
	Answer   string `json:"answer"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type UpsertDatingProfileType struct {
	AuthId *string `json:"authId"`

	Name       *string `json:"name"`
	Age        *int    `json:"age"`
	Gender     *string `json:"gender"`
	HereFor    *string `json:"hereFor"`
	LookingFor *string `json:"lookingFor"`

	Bio     *string             `json:"bio"`
	Prompts *[]DatingPromptType `json:"prompts"`

	Media *[]MediaElementType `json:"media"`

	Location      *Location `json:"location"`
	LocationLabel *string   `json:"locationLabel"`

	PreferredMatchDistance *int `json:"preferredMatchDistance"`
}

type GetPromptsType struct {
	Category *string `json:"category"`
	Page     *int64  `json:"page"`
}

type GetGendersType struct {
	Page *int64 `json:"page"`
}

type GetPromptsResponse struct {
	Records []models.Prompt `json:"records"`
	Page    *int64          `json:"page"`
	Limit   *int            `json:"limit"`
	Total   *int64          `json:"total"`
}

type GetGendersResponseType struct {
	Records []models.Gender `json:"records"`
	Page    *int64          `json:"page"`
	Limit   *int            `json:"limit"`
	Total   *int64          `json:"total"`
}

type GetProfilesType struct {
	AuthId   string `json:"authId"`
	Category string `json:"category"`
	Page     *int64 `json:"page"`
}

type GenerateMediaUploadSignedUrlType struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	AuthId      string `json:"authId"`
	FileSize    int64  `json:"fileSize"`
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
	FileSize    int64  `json:"fileSize"`
	Purpose     string `json:"purpose"`
	PartsCount  int    `json:"partsCount"`
}

type GenerateMultipartUploadUrlsResType struct {
	SignedUrls map[int]string `json:"signedUrls"`
	Expiry     int64          `json:"expiry"`
	UploadID   string         `json:"uploadID"`
	FilePath   string         `json:"filePath"`
}

type CompleteMultipartUploadType struct {
	UploadID string         `json:"uploadID"`
	FilePath string         `json:"filePath"`
	Parts    map[int]string `json:"parts"`
}

type CompleteMultipartUploadResType struct {
	URL string `json:"url"`
	ID  string `json:"id"`
}

type MediaType struct {
	ID      string `json:"_id,omitempty"`
	MediaID string `json:"id,omitempty"`
	Order   int    `json:"order,omitempty"`
	URL     string `json:"url"`
	Label   string `json:"label,omitempty"`
}

type PromptElementType struct {
	ID     string `json:"_id,omitempty"`
	Prompt string `json:"id"`
	Answer string `json:"answer"`
}

type HydratedProfileType struct {
	ID     string `json:"id"`
	AuthId string `json:"authId,omitempty"`

	ProfileCompletionScore int `json:"profileCompletionScore,omitempty"`

	Name       string `json:"name,omitempty"`
	Age        int    `json:"age,omitempty"`
	Gender     string `json:"gender,omitempty"`
	HereFor    string `json:"hereFor,omitempty"`
	LookingFor string `json:"lookingFor,omitempty"`

	Media []MediaType `json:"media,omitempty"`

	Category      string    `default:"date" json:"category,omitempty"`
	LocationLabel string    `json:"locationLabel,omitempty"`
	Location      *Location `json:"location,omitempty"`
	GeoHash       string    `json:"geoHash,omitempty"`

	Status string `default:"active" json:"status,omitempty"`

	Bio string `json:"bio,omitempty"`

	Prompts []PromptElementType `json:"prompts,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`

	PreferredMatchDistance int `bson:"preferredMatchDistance,omitempty" json:"preferredMatchDistance,omitempty"`
}
