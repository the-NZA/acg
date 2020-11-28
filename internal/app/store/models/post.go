package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post is a structure for each posts
type Post struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Excerpt  string             `bson:"excerpt,omitempty" json:"excerpt,omitempty"`
	URL      string             `bson:"url,omitempty" json:"url,omitempty"`
	Category string             `bson:"category,omitempty" json:"category,omitempty"`
	Time     time.Time          `bson:"time,omitempty" json:"time,omitempty"`
	MetaDesc string             `bson:"metadesc,omitempty" json:"metadesc,omitempty"`
	PageData []Block            `bson:"pagedata,omitempty" json:"pagedata,omitempty"`
	PostImg  string             `bson:"postimg,omitempty" json:"postimg,omitempty"`
}

// Category represents structure for each post category
type Category struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Subtitle string             `bson:"subtitle,omitempty" json:"subtitle,omitempty"`
	URL      string             `bson:"url,omitempty" json:"url,omitempty"`
	MetaDesc string             `bson:"metadesc,omitempty" json:"metadesc,omitempty"`
	Posts    []Post             `bson:"posts,omitempty" json:"posts,omitempty"`
}

// t := time.Now()
// fmt.Println(t.Format("02.01.2006"))
