package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Prompt struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Label    string `bson:"label,omitempty" json:"label,omitempty"`
	Category string `bson:"category,omitempty" json:"category,omitempty"`
	Order    int    `bson:"order,omitempty" json:"order,omitempty"`

	CreatedAt time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	DeletedAt time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}
