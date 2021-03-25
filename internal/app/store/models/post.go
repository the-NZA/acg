package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post is a structure for each posts
type Post struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Excerpt     string             `bson:"excerpt,omitempty" json:"excerpt,omitempty"`
	URL         string             `bson:"url,omitempty" json:"url,omitempty"`
	Category    string             `bson:"category,omitempty" json:"category,omitempty"`
	CategoryURL string             `bson:"categoryurl,omitempty" json:"categoryurl,omitempty"`
	Time        time.Time          `bson:"time,omitempty" json:"time,omitempty"`
	TimeString  string             `bson:"timestring,omitempty" json:"timestring,omitempty"`
	MetaDesc    string             `bson:"metadesc,omitempty" json:"metadesc,omitempty"`
	PageData    []Block            `bson:"pagedata,omitempty" json:"pagedata,omitempty"`
	PostImg     string             `bson:"postimg,omitempty" json:"postimg,omitempty"`
	Deleted     bool               `bson:"deleted" json:"-"`
}

// Category represents structure for each post category
type Category struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Subtitle string             `bson:"subtitle,omitempty" json:"subtitle,omitempty"`
	URL      string             `bson:"url,omitempty" json:"url,omitempty"`
	MetaDesc string             `bson:"metadesc,omitempty" json:"metadesc,omitempty"`
	// Posts    []Post             `bson:"posts,omitempty" json:"posts,omitempty"`
}

// t := time.Now()
// fmt.Println(t.Format("02.01.2006"))

func (p Post) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Title, validation.Required, validation.Length(5, 150)),
		validation.Field(&p.Excerpt, validation.Required, validation.Length(120, 240)),
		validation.Field(&p.URL, validation.Required, validation.Length(10, 120)),
		validation.Field(&p.Category, validation.Required, validation.Length(4, 120)),
		validation.Field(&p.CategoryURL, validation.Required, validation.Length(10, 120)),
		validation.Field(&p.MetaDesc, validation.Required, validation.Length(100, 255)),
		// validation.Field(&p.PostImg, validation.Required),
		// validation.Field(&p.PageData, validation.Required),
	)
}

func (c Category) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Title, validation.Required, validation.Length(5, 120)),
		validation.Field(&c.Subtitle, validation.Required, validation.Length(70, 255)),
		validation.Field(&c.URL, validation.Required, validation.Length(10, 120)),
		validation.Field(&c.MetaDesc, validation.Required, validation.Length(100, 255)),
	)
}
