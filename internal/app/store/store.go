package store

import (
	"context"
	"fmt"
	"time"

	"github.com/the-NZA/acg/internal/app/store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Because of one db is used decided to store name as const string
const dbName = "acg_db"

// Store – main high level abstraction on db, includes DatabaseURL string and configured *mongo.Client
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

	return nil
}

// Close just close the connection
func (s *Store) Close() {
	s.db.Disconnect(dbctx)
}

// FindOne finds all record in collection by filter
func (s *Store) FindOne(collection string, filter interface{}, opts ...*options.FindOneOptions) (bson.D, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.db.Database(dbName)
	col := db.Collection(collection)
	item := col.FindOne(ctx, filter, opts...)

	var res bson.D
	err := item.Decode(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindAll finds all items from collection
// ! Maybe deleted later
func (s *Store) FindAll(collection string) ([]bson.D, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.db.Database(dbName)
	col := db.Collection(collection)
	cur, err := col.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	items := make([]bson.D, 0)

	cur.All(ctx, &items)

	return items, nil
}

// FindAllPages returns slice of pages and error
func (s *Store) FindAllPages() ([]models.Page, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.db.Database(dbName)
	col := db.Collection("pages")
	cur, err := col.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	pages := make([]models.Page, 0)

	err = cur.All(ctx, &pages)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

// FindAllServices returns slice of pages and error
func (s *Store) FindAllServices() ([]models.Service, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.db.Database(dbName)
	col := db.Collection("services")
	cur, err := col.Find(ctx, bson.M{"deleted": false})

	if err != nil {
		return nil, err
	}

	services := make([]models.Service, 0)

	err = cur.All(ctx, &services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

// FindAllCategories returns slice of pages and error
func (s *Store) FindAllCategories(filter bson.M) ([]models.Category, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.db.Database(dbName)
	col := db.Collection("categories")
	cur, err := col.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	cats := make([]models.Category, 0)

	err = cur.All(ctx, &cats)
	if err != nil {
		return nil, err
	}

	return cats, nil
}

// FindAllPosts returns slice of pages and error
func (s *Store) FindAllPosts(filter bson.M, opts ...*options.FindOptions) ([]models.Post, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.db.Database(dbName)
	col := db.Collection("posts")
	cur, err := col.Find(ctx, filter, opts...)
	// cur, err := col.Find(ctx, bson.M{"categoryurl": "/category/news", "deleted": false}, opts...)

	if err != nil {
		return nil, err
	}

	posts := make([]models.Post, 0)

	err = cur.All(ctx, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// findPostByCategoryURL

// CountAllPosts returns slice of pages and error
func (s *Store) CountAllPosts(opts ...*options.CountOptions) (int64, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.db.Database(dbName)
	col := db.Collection("posts")
	cnt, err := col.CountDocuments(ctx, bson.M{"deleted": false}, opts...)

	if err != nil {
		return -1, err
	}

	if err != nil {
		return -1, err
	}

	return cnt, nil
}

// FindMaterials return slice of materials and error, if something went wrong
func (s *Store) FindMaterials(filter bson.M, opts ...*options.FindOptions) ([]models.Material, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.db.Database(dbName)
	col := db.Collection("materials")
	cur, err := col.Find(ctx, filter, opts...)

	if err != nil {
		return nil, err
	}

	cats := make([]models.Material, 0)

	err = cur.All(ctx, &cats)
	if err != nil {
		return nil, err
	}

	return cats, nil
}

// FindMaterials return slice of materials and error, if something went wrong
func (s *Store) FindMatcategories(filter bson.M, opts ...*options.FindOptions) ([]models.MatCategory, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.db.Database(dbName)
	col := db.Collection("matcategories")
	cur, err := col.Find(ctx, filter, opts...)

	if err != nil {
		return nil, err
	}

	cats := make([]models.MatCategory, 0)

	err = cur.All(ctx, &cats)
	if err != nil {
		return nil, err
	}

	return cats, nil
}

// InsertOne add data to collection
func (s *Store) InsertOne(collection string, data interface{}) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.db.Database(dbName)
	col := db.Collection(collection)
	r, err := col.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}
	fmt.Println("Inserted to: ", col.Name())

	return r, nil
}

// UpdateOne updates one record
func (s *Store) UpdateOne(collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.db.Database(dbName)
	col := db.Collection(collection)

	r, err := col.UpdateOne(
		ctx,
		filter,
		update,
	)

	if err != nil {
		return nil, err
	}

	return r, nil
}

// func (s *Store) DeleteOne(collection string, filter interface{}) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	db := s.db.Database(dbName)
// 	col := db.Collection(collection)

// 	r, err := col.DeleteOne(ctx, filter, )

// 	if err != nil {
// 		return nil, err
// 	}

// 	return r, nil
// }
