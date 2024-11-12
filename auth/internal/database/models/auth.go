package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth struct {
	ID                       primitive.ObjectID `bson:"_id,omitempty"`
	Phone                    string             `bson:"phone" unique:"true"`
	PhoneConfirmedAtUtcEpoch int8               `bson:"phoneConfirmedAtUtcEpoch,omitempty"`
	Status                   string             `bson:"status,omitempty"`
	CreatedAt                time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt                time.Time          `bson:"updatedAt,omitempty"`
}
