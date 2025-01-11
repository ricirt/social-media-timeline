package model

import (
	"time"

	pkgErrors "github.com/ricirt/social-media-timeline/services/pkg/errors"
)

// User represents the user model in the application
type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	CreatedAt time.Time `json:"created_at"`
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
