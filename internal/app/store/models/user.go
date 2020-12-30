package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represets each user
type User struct {
	ID                primitive.ObjectID `bson:"_id" json:"_id"`
	Username          string             `bson:"username" json:"username"`
	Password          string             `bson:"-" json:"-"`
	EncryptedPassword string             `bson:"pswd" json:"pswd"`
	Email             string             `bson:"email" json:"email"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required, validation.Length(10, 30)),
		validation.Field(&u.EncryptedPassword, validation.Required),
		validation.Field(&u.Email, is.Email),
	)
}
