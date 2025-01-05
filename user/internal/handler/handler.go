package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ricirt/social-media-timeline/user/internal/repository/mongo"
)

// UserHandler defines CRUD methods
type UserHandler interface {
	GetUsers(c *gin.Context)
	GetUserByID(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

// userHandler struct
type userHandler struct {
	repo *mongo.MongoUserRepository
}

// NewUserHandler returns a new instance of userHandler
func NewUserHandler(repo *mongo.MongoUserRepository) UserHandler {
	return &userHandler{
		repo: repo,
	}
}

// GetUsers returns all users
func (h *userHandler) GetUsers(c *gin.Context) {
	users, err := h.repo.GetUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID returns a user by ID
func (h *userHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.repo.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser creates a new user
func (h *userHandler) CreateUser(c *gin.Context) {
	var newUser mongo.User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.CreateUser(c, &newUser)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newUser)
}

// UpdateUser updates an existing user
func (h *userHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser mongo.User

	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateUser(c, id, &updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser implements UserHandler
func (h *userHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteUser(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
