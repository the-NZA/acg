package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Page is basic model for each page
type Page struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title" json:"title"`
	MetaDesc string             `bson:"desc,omitempty" json:"desc,omitempty"`
	Slug     string             `bson:"slug,omitempty" json:"slug,omitempty"`
}

func (p *Page) GetAll() ([]Page, error) {
	return nil, nil
}
