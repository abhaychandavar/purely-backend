package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	URL         string             `bson:"url,omitempty" json:"url" unique:"true"`
	EXT         string             `bson:"ext,omitempty" json:"ext"`
	RefID       primitive.ObjectID `bson:"refID,omitempty" json:"refID" unique:"true"`
	Domain      string             `bson:"domain,omitempty" json:"domain"`
	Path        string             `bson:"path,omitempty" json:"path"`
	ContentType string             `bson:"contentType,omitempty" json:"contentType"`
	FileName    string             `bson:"fileName,omitempty" json:"fileName"`
	Size        int                `bson:"size,omitempty" json:"size"`
}
