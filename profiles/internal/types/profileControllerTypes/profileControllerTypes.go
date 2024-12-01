package profileControllerTypes

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateProfileType struct {
	Lat *float64 `json:"lat"`
	Lng *float64 `json:"lng"`
}

type ImageElementType struct {
	ImageId *primitive.ObjectID `json:"id"`
	Order   *int                `json:"order"`
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
	Images                 *[]ImageElementType `json:"images"`
	Lat                    *float64            `json:"lat"`
	Lng                    *float64            `json:"lng"`
	LocationLabel          *string             `json:"locationLabel"`
	Location               *Location           `json:"location"`
	PreferredMatchDistance *int                `json:"preferredMatchDistance"`
}
