package mongodb

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"

	"GoNews/pkg/storage"
)

type Store struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func NewStorage(connectionString, dbName, collectionName string) (*Store, error) {
	client, err := mongo.Connect(context.Background(), connectionString, nil)

	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)
	collection := database.Collection(collectionName)

	return &Store{client, database, collection}, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []storage.Post
	for cursor.Next(ctx) {
		var p storage.Post
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Store) AddPost(post storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.collection.InsertOne(ctx, post)
	return err
}

func (s *Store) UpdatePost(post storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": post.ID}
	update := bson.M{
		"$set": bson.M{
			"title":        post.Title,
			"content":      post.Content,
			"author_id":    post.AuthorID,
			"created_at":   post.CreatedAt,
			"published_at": post.PublishedAt,
			"author_name":  post.AuthorName,
		},
	}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *Store) DeletePost(post storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": post.ID}
	_, err := s.collection.DeleteOne(ctx, filter)
	return err
}
