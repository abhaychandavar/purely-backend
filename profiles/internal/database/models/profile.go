package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Type        string    `bson:"type,omitempty" default:"Point" json:"type,omitempty"`
	Coordinates []float64 `bson:"coordinates,omitempty" json:"coordinates,omitempty"`
}

type ImageElementType struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ImageId         primitive.ObjectID `bson:"image,omitempty" json:"id,omitempty"`
	Order           int                `bson:"order,omitempty" json:"order,omitempty"`
	ImageUrl        string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	BlurredImageUrl string             `bson:"blurredImageUrl,omitempty" json:"blurredImageUrl,omitempty"`
}

type PromptElementType struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Prompt primitive.ObjectID `bson:"prompt,omitempty" json:"id"`
	Answer string             `bson:"answer,omitempty"  json:"answer"`
}

type Profile struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AuthId string             `bson:"authId" json:"authId,omitempty"`

	ProfileCompletionScore int `bson:"profileCompletionScore,omitempty" json:"profileCompletionScore,omitempty"`

	Name       string             `bson:"name,omitempty" json:"name,omitempty"`
	Age        int                `bson:"age,omitempty" json:"age,omitempty"`
	Gender     primitive.ObjectID `bson:"gender,omitempty" json:"gender,omitempty"`
	HereFor    string             `bson:"hereFor,omitempty" json:"hereFor,omitempty"`
	LookingFor string             `bson:"lookingFor,omitempty" json:"lookingFor,omitempty"`

	Images []ImageElementType `bson:"images,omitempty" json:"images,omitempty"`

	Category      string    `bson:"category,omitempty" default:"date" json:"category,omitempty"`
	LocationLabel string    `bson:"locationLabel,omitempty" json:"locationLabel,omitempty"`
	Location      *Location `bson:"location,omitempty" json:"location,omitempty"`
	GeoHash       string    `bson:"geohash,omitempty" json:"geoHash,omitempty"`

	Status string `bson:"status,omitempty" default:"active" json:"status,omitempty"`

	Bio string `bson:"bio,omitempty" json:"bio,omitempty"`

	Prompts []PromptElementType `bson:"prompts,omitempty" json:"prompts,omitempty"`

	CreatedAt time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	DeletedAt time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`

	PreferredMatchDistance int `bson:"preferredMatchDistance,omitempty" json:"preferredMatchDistance,omitempty"`
}
