package mongo

import (
	"context"
	"fmt"

	"github.com/ricirt/social-media-timeline/user/internal/model"
	pkgErrors "github.com/ricirt/social-media-timeline/user/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUsers(ctx context.Context) ([]model.User, error)
	UpdateUser(ctx context.Context, id string, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
}

// MongoUserRepository implements UserRepository interface using MongoDB
type MongoUserRepository struct {
	collection     *mongo.Collection
	DbName         string
	CollectionName string
}

// NewMongoUserRepository creates and initializes a new MongoUserRepository
// It also sets up necessary indexes for the collection
func NewMongoUserRepository(client *mongo.Client, dbName, collectionName string) (*MongoUserRepository, error) {
	if client == nil {
		return nil, pkgErrors.ErrDBConnection
	}

	collection := client.Database(dbName).Collection(collectionName)

	// Create a unique index for the email field
	_, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

	return &MongoUserRepository{
		collection:     collection,
		DbName:         dbName,
		CollectionName: collectionName,
	}, nil
}

// CreateUser creates a new user in the database
func (r *MongoUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return pkgErrors.ErrDuplicateEmail
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return nil
}

// GetUserByID retrieves a user by their ID
func (r *MongoUserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, pkgErrors.ErrInvalidID
	}

	var user model.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, pkgErrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUsers retrieves all users from the database
func (r *MongoUserRepository) GetUsers(ctx context.Context) ([]model.User, error) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

// UpdateUser updates an existing user's information
func (r *MongoUserRepository) UpdateUser(ctx context.Context, id string, user *model.User) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return pkgErrors.ErrInvalidID
	}

	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return pkgErrors.ErrDuplicateEmail
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return pkgErrors.ErrUserNotFound
	}

	return nil
}

// DeleteUser removes a user from the database by their ID
func (r *MongoUserRepository) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return pkgErrors.ErrInvalidID
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return pkgErrors.ErrUserNotFound
	}

	return nil
}
