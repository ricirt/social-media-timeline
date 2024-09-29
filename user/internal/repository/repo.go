package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User struct'ı: MongoDB'de tutulacak kullanıcı modeli
type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string             `bson:"name" json:"name"`
	Email string             `bson:"email" json:"email"`
}

// UserRepository defines the interface for user repository
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, id string, user *User) error
	DeleteUser(ctx context.Context, id string) error
}

// MongoUserRepository implements UserRepository interface
type MongoUserRepository struct {
	collection     *mongo.Collection
	DbName         string
	CollectionName string
}

// NewMongoUserRepository creates a new MongoUserRepository
func NewMongoUserRepository(client *mongo.Client, dbName, collectionName string) *MongoUserRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &MongoUserRepository{collection: collection, DbName: dbName, CollectionName: collectionName}
}

// CreateUser inserts a new user into the collection
func (r *MongoUserRepository) CreateUser(ctx context.Context, user *User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// GetUserByID retrieves a user by ID
func (r *MongoUserRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return &user, err
}

// GetUserByID retrieves a user by ID
func (r *MongoUserRepository) GetUsers(ctx context.Context) ([]User, error) {
	return nil, nil
}

// UpdateUser updates an existing user
func (r *MongoUserRepository) UpdateUser(ctx context.Context, id string, user *User) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

// DeleteUser deletes a user by ID
func (r *MongoUserRepository) DeleteUser(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
