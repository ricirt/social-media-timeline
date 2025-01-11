package model

import (
	pkgErrors "github.com/ricirt/social-media-timeline/services/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the user model in the application
type User struct {
	ID    interface{} `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string      `json:"name" validate:"required"`
	Email string      `json:"email" validate:"required,email"`
}

// MongoUser converts User to MongoDB specific model
func (u *User) MongoUser() *User {
	if u.ID != nil {
		if id, ok := u.ID.(string); ok {
			if objectID, err := primitive.ObjectIDFromHex(id); err == nil {
				u.ID = objectID
			}
		}
	}
	return u
}

// PostgresUser converts User to PostgreSQL specific model
func (u *User) PostgresUser() *User {
	if u.ID != nil {
		if objectID, ok := u.ID.(primitive.ObjectID); ok {
			u.ID = objectID.Hex()
		}
	}
	return u
}

// Validate performs validation checks on user data
func (u *User) Validate() error {
	if u.Name == "" {
		return pkgErrors.ErrEmptyName
	}
	if u.Email == "" {
		return pkgErrors.ErrEmptyEmail
	}
	return nil
}
