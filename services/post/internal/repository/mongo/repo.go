package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	pkgErrors "github.com/ricirt/social-media-timeline/services/pkg/errors"
	"github.com/ricirt/social-media-timeline/services/post/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoPostRepository struct {
	collection *mongo.Collection
}

func NewMongoPostRepository(client *mongo.Client, dbName, collectionName string) (*MongoPostRepository, error) {
	if client == nil {
		return nil, pkgErrors.ErrDBConnection
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoPostRepository{collection: collection}, nil
}

func (r *MongoPostRepository) CreatePost(ctx context.Context, post *model.Post) error {
	post.ID = uuid.New().String()
	post.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, post)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

func (r *MongoPostRepository) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	var post model.Post
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	if err == mongo.ErrNoDocuments {
		return nil, pkgErrors.ErrPostNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return &post, nil
}

func (r *MongoPostRepository) GetPosts(ctx context.Context) ([]model.Post, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer cursor.Close(ctx)

	var posts []model.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("failed to decode posts: %w", err)
	}
	return posts, nil
}

func (r *MongoPostRepository) GetPostsByUserID(ctx context.Context, userID string) ([]model.Post, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer cursor.Close(ctx)

	var posts []model.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("failed to decode posts: %w", err)
	}
	return posts, nil
}

func (r *MongoPostRepository) UpdatePost(ctx context.Context, id string, post *model.Post) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return pkgErrors.ErrInvalidID
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"user_id": post.UserID,
			"content": post.Content,
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	if result.MatchedCount == 0 {
		return pkgErrors.ErrPostNotFound
	}
	return nil
}

func (r *MongoPostRepository) DeletePost(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return pkgErrors.ErrInvalidID
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	if result.DeletedCount == 0 {
		return pkgErrors.ErrPostNotFound
	}
	return nil
}
