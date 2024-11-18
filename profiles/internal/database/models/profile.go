package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Type        string    `bson:"type" default:"Point"`
	Coordinates []float64 `bson:"coordinates"`
}

type Profile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Category  string             `bson:"category" default:"date"`
	AuthId    string             `bson:"authId" unique:"true"`
	Location  Location           `bson:"location"`
	GeoHash   string             `bson:"geohash"`
	Status    string             `bson:"status" default:"active"`
	CreatedAt time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty"`
}
