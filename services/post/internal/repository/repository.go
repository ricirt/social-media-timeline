package repository

import (
	"context"
	"database/sql"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	pkgErrors "github.com/ricirt/social-media-timeline/services/pkg/errors"
	"github.com/ricirt/social-media-timeline/services/post/internal/model"
	mongoRepo "github.com/ricirt/social-media-timeline/services/post/internal/repository/mongo"
	pgRepo "github.com/ricirt/social-media-timeline/services/post/internal/repository/postgres"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
	GetPosts(ctx context.Context) ([]model.Post, error)
	GetPostsByUserID(ctx context.Context, userID string) ([]model.Post, error)
	UpdatePost(ctx context.Context, id string, post *model.Post) error
	DeletePost(ctx context.Context, id string) error
}

func NewMongoPostRepository(client interface{}) (PostRepository, error) {
	mongoClient, ok := client.(*mongo.Client)
	if !ok {
		return nil, errors.New("invalid mongodb client")
	}
	return mongoRepo.NewMongoPostRepository(mongoClient, "timeline", "posts")
}

func NewPostgresPostRepository(client interface{}) (PostRepository, error) {
	pgClient, ok := client.(*sql.DB)
	if !ok {
		return nil, errors.New("invalid postgres client")
	}
	return pgRepo.NewPostgresPostRepository(pgClient)
}

func NewPostRepository(dbType string, dbClient interface{}) (PostRepository, error) {
	switch dbType {
	case "mongodb":
		return NewMongoPostRepository(dbClient)
	case "postgres":
		return NewPostgresPostRepository(dbClient)
	default:
		return nil, pkgErrors.ErrUnsupportedDatabase
	}
}
