package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	Url string `bson:"url" json:"url" unique:"true"`
	EXT string `bson:"ext" json:"ext"`
}
