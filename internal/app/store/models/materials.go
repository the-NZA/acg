package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Material represent structure of each material
type Material struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Title      string             `bson:"title,omitempty" json:"title,omitempty"`
	Category   string             `bson:"category,omitempty" json:"category,omitempty"`
	Desc       string             `bson:"desc,omitempty" json:"desc,omitempty"`
	Time       time.Time          `bson:"time,omitempty" json:"time,omitempty"`
	TimeString string             `bson:"timestring,omitempty" json:"timestring,omitempty"`
	FileLink   string             `bson:"filelink,omitempty" json:"filelink,omitempty"`
}

// MatCategory represent each materials category
type MatCategory struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Title     string             `bson:"title,omitempty" json:"title,omitempty"`
	Slug      string             `bson:"slug,omitempty" json:"slug,omitempty"`
	Desc      string             `bson:"desc,omitempty" json:"desc,omitempty"`
	Materials []Material         `bson:"materials,omitempty" json:"materials,omitempty"`
}

func (m Material) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Title, validation.Required, validation.Length(5, 150)),
		validation.Field(&m.Category, validation.Required, validation.Length(5, 150)),
		validation.Field(&m.Desc, validation.Required, validation.Length(70, 255)),
		validation.Field(&m.FileLink, validation.Required),
	)
}

func (m MatCategory) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Title, validation.Required, validation.Length(5, 150)),
		validation.Field(&m.Desc, validation.Required, validation.Length(70, 255)),
	)
}
