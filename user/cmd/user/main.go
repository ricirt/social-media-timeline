package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/ricirt/social-media-timeline/user/internal/handler"
)

// Kullanıcı yapısı
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Fake kullanıcı verisi (örnek olarak hafızada tutuyoruz)
var users = []User{
	{ID: "1", Name: "John Doe", Email: "john@example.com"},
	{ID: "2", Name: "Jane Smith", Email: "jane@example.com"},
}

func main() {

	// Gin router'ı başlat
	r := gin.Default()

	userHandler := handler.NewUserHandler()
	// Kullanıcı yolları için route'ları tanımla
	r.GET("/users", getUsers)           // Tüm kullanıcıları listele
	r.GET("/users/:id", getUserByID)    // ID'ye göre kullanıcı getir
	r.POST("/users", createUser)        // Yeni kullanıcı yarat
	r.PUT("/users/:id", updateUser)     // Mevcut kullanıcıyı güncelle
	r.DELETE("/users/:id", deleteUser)  // Kullanıcıyı sil

	// Sunucuyu başlat
	r.Run(":8080")  // 8080 portunda çalıştır
}