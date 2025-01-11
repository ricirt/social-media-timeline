package model

import (
	"time"

	pkgErrors "github.com/ricirt/social-media-timeline/services/pkg/errors"
)

type Post struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Post) Validate() error {
	if p.UserID == "" {
		return pkgErrors.ErrEmptyUserID
	}
	if p.Content == "" {
		return pkgErrors.ErrEmptyContent
	}
	return nil
}
