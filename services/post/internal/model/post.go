package model

import (
	pkgErrors "github.com/ricirt/social-media-timeline/services/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID      interface{} `json:"id,omitempty" bson:"_id,omitempty"`
	UserID  string      `json:"user_id" validate:"required"`
	Content string      `json:"content" validate:"required"`
}

func (p *Post) MongoPost() *Post {
	if p.ID != nil {
		if id, ok := p.ID.(string); ok {
			if objectID, err := primitive.ObjectIDFromHex(id); err == nil {
				p.ID = objectID
			}
		}
	}
	return p
}

func (p *Post) PostgresPost() *Post {
	if p.ID != nil {
		if objectID, ok := p.ID.(primitive.ObjectID); ok {
			p.ID = objectID.Hex()
		}
	}
	return p
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
