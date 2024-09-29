package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"

)

// UserHandler interface'i CRUD metotlarını tanımlar
type UserHandler interface {
	GetUsers(c *gin.Context)
	GetUserByID(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

// User struct'ı
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Mock kullanıcı verisi
var users = []User{
	{ID: "1", Name: "John Doe", Email: "john@example.com"},
	{ID: "2", Name: "Jane Smith", Email: "jane@example.com"},
}

// userHandler yapısı
type userHandler struct{}

// NewUserHandler, userHandler'ın bir örneğini döner
func NewUserHandler() UserHandler {
	return &userHandler{}
}

// GetUsers tüm kullanıcıları döndürür
func (h *userHandler) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

// GetUserByID, ID'ye göre bir kullanıcı döner
func (h *userHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

// CreateUser yeni bir kullanıcı oluşturur
func (h *userHandler) CreateUser(c *gin.Context) {
	var newUser User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users = append(users, newUser)
	c.JSON(http.StatusCreated, newUser)
}

// UpdateUser, mevcut kullanıcıyı günceller
func (h *userHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser User

	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, user := range users {
		if user.ID == id {
			users[i] = updatedUser
			c.JSON(http.StatusOK, updatedUser)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

// DeleteUser, kullanıcıyı siler
func (h *userHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H