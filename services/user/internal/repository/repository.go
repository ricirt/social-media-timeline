package repository

import (
	"context"
	"database/sql"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ricirt/social-media-timeline/services/user/internal/model"
	mongoRepo "github.com/ricirt/social-media-timeline/services/user/internal/repository/mongo"
	pgRepo "github.com/ricirt/social-media-timeline/services/user/internal/repository/postgres"
	pkgErrors "github.com/ricirt/social-media-timeline/services/user/pkg/errors"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUsers(ctx context.Context) ([]model.User, error)
	UpdateUser(ctx context.Context, id string, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
}

func NewMongoUserRepository(client interface{}) (UserRepository, error) {
	mongoClient, ok := client.(*mongo.Client)
	if !ok {
		return nil, errors.New("invalid mongodb client")
	}
	return mongoRepo.NewMongoUserRepository(mongoClient, "timeline", "users")
}

func NewPostgresUserRepository(client interface{}) (UserRepository, error) {
	pgClient, ok := client.(*sql.DB)
	if !ok {
		return nil, errors.New("invalid postgres client")
	}
	return pgRepo.NewPostgresUserRepository(pgClient)
}

// NewUserRepository creates a new user repository based on the database type
func NewUserRepository(dbType string, dbClient interface{}) (UserRepository, error) {
	switch dbType {
	case "mongodb":
		return NewMongoUserRepository(dbClient)
	case "postgres":
		return NewPostgresUserRepository(dbClient)
	default:
		return nil, pkgErrors.ErrUnsupportedDatabase
	}
}
