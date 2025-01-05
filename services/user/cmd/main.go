package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ricirt/social-media-timeline/services/user/internal/handler"
	"github.com/ricirt/social-media-timeline/services/user/internal/repository/mongo"
	"github.com/ricirt/social-media-timeline/services/user/pkg/config"
	mongodb "github.com/ricirt/social-media-timeline/services/user/pkg/mongo-db"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	if err := mongodb.InitMongoClient(config.MongoDB.ConnectionString); err != nil {
		log.Fatal("Failed to initialize MongoDB client:", err)
	}

	mongoClient, err := mongodb.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	mongoUserRepository, err := mongo.NewMongoUserRepository(
		mongoClient,
		config.MongoDB.Database,
		config.MongoDB.Collections.Users,
	)
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewUserHandler(mongoUserRepository)
	r := gin.Default()

	// Define routes for user endpoints
	r.GET("/users", h.GetUsers)          // List all users
	r.GET("/users/:id", h.GetUserByID)   // Get user by ID
	r.POST("/users", h.CreateUser)       // Create new user
	r.PUT("/users/:id", h.UpdateUser)    // Update existing user
	r.DELETE("/users/:id", h.DeleteUser) // Delete user

	// Start the server
	r.Run(":" + config.Server.Port)
}
