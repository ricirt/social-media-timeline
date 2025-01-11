package errors

import "errors"

var (
    // Database errors
    ErrUnsupportedDatabase = errors.New("unsupported database type")
    ErrDBConnection       = errors.New("database connection error")
    
    // Validation errors
    ErrEmptyUserID       = errors.New("user id cannot be empty")
    ErrEmptyContent      = errors.New("content cannot be empty")
    ErrEmptyName         = errors.New("name cannot be empty")
    ErrEmptyEmail        = errors.New("email cannot be empty")
    
    // Not found errors
    ErrUserNotFound      = errors.New("user not found")
    ErrPostNotFound      = errors.New("post not found")
    
    // Other errors
    ErrInvalidID         = errors.New("invalid id")
    ErrDuplicateEmail    = errors.New("email already exists")
)
