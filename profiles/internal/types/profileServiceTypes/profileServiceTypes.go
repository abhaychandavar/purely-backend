package profileServiceTypes

import (
	"profiles/internal/database/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

type ImageElementType struct {
	ImageId primitive.ObjectID `json:"id"`
	Order   int                `json:"order"`
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

	Images *[]ImageElementType `json:"images"`

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
	FileName string `json:"fileName"`
	MimeType string `json:"mimeType"`
	AuthId   string `json:"authId"`
	FileSize int64  `json:"fileSize"`
	Purpose  string `json:"purpose"`
}

type GenerateMediaUploadSignedUrlResType struct {
	SignedUrl string `json:"signedUrl"`
	Expiry    int64  `json:"expiry"`
}
