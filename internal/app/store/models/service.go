package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service is a structure for representing each service
type Service struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Img      *ServiceImage      `bson:"img,omitempty" json:"img,omitempty"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Subtitle string             `bson:"subtitle,omitempty" json:"subtitle,omitempty"`
	Desc     string             `bson:"desc,omitempty" json:"desc,omitempty"`
	Deleted  bool               `bson:"deleted" json:"-"`
}

// ServiceImage represets basic structure of service card image
type ServiceImage struct {
	URL string `bson:"url,omitempty" json:"url,omitempty"`
	Alt string `bson:"alt,omitempty" json:"alt,omitempty"`
}

func (s Service) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Img, validation.Required),
		validation.Field(&s.Title, validation.Required, validation.Length(5, 150)),
		validation.Field(&s.Subtitle, validation.Required, validation.Length(70, 255)),
		validation.Field(&s.Desc, validation.Required, validation.Length(80, 255)),
	)
}

// TODO: Decide add or not validation for ServiceImage
