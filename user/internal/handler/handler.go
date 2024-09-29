package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ricirt/social-media-timeline/user/internal/repository"
)

// UserHandler interface'i CRUD metotlarını tanımlar
type UserHandler interface {
	GetUsers(c *gin.Context)
	GetUserByID(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

// userHandler yapısı
type userHandler struct {
	repo *repository.MongoUserRepository
}

// NewUserHandler, userHandler'ın bir örneğini döner
func NewUserHandler(repo *repository.MongoUserRepository) UserHandler {
	return &userHandler{
		repo: repo,
	}
}

// GetUsers tüm kullanıcıları döndürür
func (h *userHandler) GetUsers(c *gin.Context) {

	users, err := h.repo.GetUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID, ID'ye göre bir kullanıcı döner
func (h *userHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.repo.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})

		return
	}
	c.JSON(http.StatusOK, user)
	return
}

// CreateUser yeni bir kullanıcı oluşturur
func (h *userHandler) CreateUser(c *gin.Context) {
	var newUser repository.User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.CreateUser(c, &newUser)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, nil)
}

// UpdateUser, mevcut kullanıcıyı günceller
func (h *userHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser repository.User

	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateUser(c, id, &updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
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
