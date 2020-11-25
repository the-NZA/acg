package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName = "acg_db"

// Store – main high level abstraction on db
type Store struct {
	DatabaseURL string
	db          *mongo.Client
	// DB *mongo.Database
}

// New creates new store
func New(dburl string) *Store {
	return &Store{
		DatabaseURL: dburl,
	}
}

// context for mongo db
var dbctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

// Open just open new connection
func (s *Store) Open() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(s.DatabaseURL))
	if err != nil {
		return err
	}
	err = client.Connect(dbctx)
	if err != nil {
		return err
	}

	if err = client.Ping(dbctx, nil); err != nil {
		return err
	}

	s.db = client
	// s.DB = client.Database("acg_db")
	return nil
}

// Close just close the connection
func (s *Store) Close() {
	// s.DB.Client().Disconnect(dbctx)
	s.db.Disconnect(dbctx)
}

// FindAll finds all items from collection
func (s *Store) FindAll(collection string) ([]bson.D, error) {
	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	db := s.db.Database(dbName)
	col := db.Collection(collection)
	cur, err := col.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	items := []bson.D{}

	cur.All(ctx, &items)
	// i := bson.D{}

	// for cur.Next(ctx) {
	// 	cur.Decode(&i)
	// 	items = append(items, i)
	// }

	return items, nil
}

// InsertOne add data to collection
func (s *Store) InsertOne(collection string, data interface{}) (*mongo.InsertOneResult, error) {
	// col := s.DB.Collection(collection)
	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	db := s.db.Database(dbName)
	col := db.Collection(collection)
	r, err := col.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}
	fmt.Println(col.Name())

	return r, nil
}
