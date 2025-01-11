package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	pkgErrors "github.com/ricirt/social-media-timeline/services/pkg/errors"
	"github.com/ricirt/social-media-timeline/services/post/internal/model"
)

type PostgresPostRepository struct {
	db *sql.DB
}

func NewPostgresPostRepository(db *sql.DB) (*PostgresPostRepository, error) {
	if db == nil {
		return nil, pkgErrors.ErrDBConnection
	}
	return &PostgresPostRepository{db: db}, nil
}

func (r *PostgresPostRepository) CreatePost(ctx context.Context, post *model.Post) error {
	query := `
        INSERT INTO posts (user_id, content)
        VALUES ($1, $2)
        RETURNING id`

	err := r.db.QueryRowContext(ctx, query, post.UserID, post.Content).Scan(&post.ID)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

func (r *PostgresPostRepository) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, pkgErrors.ErrInvalidID
	}

	query := `
        SELECT id, user_id, content
        FROM posts
        WHERE id = $1`

	post := &model.Post{}
	err = r.db.QueryRowContext(ctx, query, idInt).Scan(&post.ID, &post.UserID, &post.Content)
	if err == sql.ErrNoRows {
		return nil, pkgErrors.ErrPostNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return post, nil
}

func (r *PostgresPostRepository) GetPosts(ctx context.Context) ([]model.Post, error) {
	query := `
        SELECT id, user_id, content
        FROM posts`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostgresPostRepository) GetPostsByUserID(ctx context.Context, userID string) ([]model.Post, error) {
	query := `
        SELECT id, user_id, content
        FROM posts
        WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostgresPostRepository) UpdatePost(ctx context.Context, id string, post *model.Post) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return pkgErrors.ErrInvalidID
	}

	query := `
        UPDATE posts
        SET user_id = $1, content = $2
        WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, post.UserID, post.Content, idInt)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return pkgErrors.ErrPostNotFound
	}
	return nil
}

func (r *PostgresPostRepository) DeletePost(ctx context.Context, id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return pkgErrors.ErrInvalidID
	}

	query := `
        DELETE FROM posts
        WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, idInt)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return pkgErrors.ErrPostNotFound
	}
	return nil
}
