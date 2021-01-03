package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User represets each user
type User struct {
	ID                primitive.ObjectID `bson:"_id" json:"_id"`
	Username          string             `bson:"username" json:"username"`
	EncryptedPassword string             `bson:"pswd" json:"-"`
	Email             string             `bson:"email,omitempty" json:"email,omitempty"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required, validation.Length(10, 30)),
		validation.Field(&u.EncryptedPassword, validation.Required),
		validation.Field(&u.Email, is.Email),
	)
}

func (u *User) HashPassword(password string) error {
	bts, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	u.EncryptedPassword = string(bts)
	return nil
}

func (u User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
	return err == nil
}
