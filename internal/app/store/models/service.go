package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Service is a structure for representing each service
type Service struct {
	ID  primitive.ObjectID `bson:"_id" json:"_id"`
	Img struct {
		URL string `bson:"url" json:"url"`
		Alt string `bson:"alt" json:"alt"`
	} `bson:"img" json:"img"`
	Title    string `bson:"title" json:"title"`
	Subtitle string `bson:"subtitle" json:"subtitle"`
	Desc     string `bson:"desc" json:"desc"`
}
