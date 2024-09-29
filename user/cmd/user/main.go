package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ricirt/social-media-timeline/user/internal/handler"
	"github.com/ricirt/social-media-timeline/user/internal/repository"
	mongodb "github.com/ricirt/social-media-timeline/user/pkg/mongo-db"
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

	mongodb.InitMongoClient("mongodb://localhost:27017")
	mongoClient := mongodb.GetMongoClient()
	mongoUserRepository := repository.NewMongoUserRepository(mongoClient, "social-media-timeline", "users")
	h := handler.NewUserHandler(mongoUserRepository)
	r := gin.Default()

	// Kullanıcı yolları için route'ları tanımla
	r.GET("/users", h.GetUsers)          // Tüm kullanıcıları listele
	r.GET("/users/:id", h.GetUserByID)   // ID'ye göre kullanıcı getir
	r.POST("/users", h.CreateUser)       // Yeni kullanıcı yarat
	r.PUT("/users/:id", h.UpdateUser)    // Mevcut kullanıcıyı güncelle
	r.DELETE("/users/:id", h.DeleteUser) // Kullanıcıyı sil

	// Sunucuyu başlat
	r.Run(":8080") // 8080 portunda çalıştır
}
