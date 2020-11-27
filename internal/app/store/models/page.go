package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Page is basic model for each page
type Page struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Subtitle string             `bson:"subtitle,omitempty" json:"subtitle,omitempty"`
	MetaDesc string             `bson:"desc,omitempty" json:"desc,omitempty"`
	Slug     string             `bson:"slug,omitempty" json:"slug,omitempty"`
	PageData []Block            `bson:"pagedata,omitempty" json:"pagedata,omitempty"`
}

// Block is a struct that represents blocks saved in json format
// Specific block types: text, image, ...
type Block struct {
	Type string     `bson:"type" json:"type"`
	Data *BlockData `bson:"data" json:"data"`
}

// BlockData represents structure of blocks data field
type BlockData struct {
	Text  string    `bson:"text,omitempty" json:"text,omitempty"`
	Level int8      `bson:"level,omitempty" json:"level,omitempty"`
	File  *FileInfo `bson:"file,omitempty" json:"file,omitempty"`
}

// FileInfo represents basic file structure for editors "file" field
type FileInfo struct {
	URL    string `bson:"url,omitempty" json:"url,omitempty"`
	Width  int    `bson:"width,omitempty" json:"width,omitempty"`
	Height int    `bson:"height,omitempty" json:"height,omitempty"`
}
