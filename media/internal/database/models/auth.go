package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Identifier string             `bson:"identifier,omitempty" unique:"true" required:"true"`
	Phone      string             `bson:"phone,omitempty" unique:"true"`
	Status     string             `bson:"status,omitempty"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt  time.Time          `bson:"updatedAt,omitempty"`
}
