package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/ricirt/social-media-timeline/services/user/internal/model"
	pkgErrors "github.com/ricirt/social-media-timeline/services/pkg/errors"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUsers(ctx context.Context) ([]model.User, error)
	UpdateUser(ctx context.Context, id string, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
}

// PostgresUserRepository implements UserRepository interface using PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgresUserRepository
func NewPostgresUserRepository(db *sql.DB) (*PostgresUserRepository, error) {
	if db == nil {
		return nil, pkgErrors.ErrDBConnection
	}

	return &PostgresUserRepository{
		db: db,
	}, nil
}

// CreateUser creates a new user
func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	user := &model.User{}
	query := `
		SELECT id, name, email
		FROM users
		WHERE id = $1`

	err = r.db.QueryRowContext(ctx, query, idInt).Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		return nil, pkgErrors.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUsers retrieves all users
func (r *PostgresUserRepository) GetUsers(ctx context.Context) ([]model.User, error) {
	query := `
		SELECT id, name, email
		FROM users
		ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}

	return users, nil
}

// UpdateUser updates an existing user
func (r *PostgresUserRepository) UpdateUser(ctx context.Context, id string, user *model.User) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	query := `
		UPDATE users
		SET name = $1, email = $2
		WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, idInt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return pkgErrors.ErrUserNotFound
	}

	return nil
}

// DeleteUser removes a user
func (r *PostgresUserRepository) DeleteUser(ctx context.Context, id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	query := `
		DELETE FROM users
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, idInt)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return pkgErrors.ErrUserNotFound
	}

	return nil
}
