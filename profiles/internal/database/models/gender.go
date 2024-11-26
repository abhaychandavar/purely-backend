package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gender struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Label       string `bson:"label,omitempty" json:"label,omitempty"`
	Code        string `bson:"code,omitempty" json:"code,omitempty"`
	Order       int    `bson:"order,omitempty" json:"order,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`

	CreatedAt time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	DeletedAt time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}
