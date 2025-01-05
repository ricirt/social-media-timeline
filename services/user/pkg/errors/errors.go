package errors

import "errors"

// Common repository errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidID    = errors.New("invalid user ID")
	ErrEmptyName    = errors.New("user name cannot be empty")
	ErrEmptyEmail   = errors.New("user email cannot be empty")
)

// Database operation errors
var (
	ErrDuplicateEmail      = errors.New("user with this email already exists")
	ErrDBConnection        = errors.New("database connection error")
	ErrDBOperation         = errors.New("database operation failed")
	ErrUnsupportedDatabase = errors.New("unsupported database type")
)

// Validation errors
var (
	ErrValidation = errors.New("validation error")
)
