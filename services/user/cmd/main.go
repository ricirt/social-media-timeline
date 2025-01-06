package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ricirt/social-media-timeline/services/user/internal/handler"
	"github.com/ricirt/social-media-timeline/services/user/internal/repository"
	"github.com/ricirt/social-media-timeline/services/user/pkg/config"

	//mongodb "github.com/ricirt/social-media-timeline/services/user/pkg/mongo-db"
	"log"

	"github.com/ricirt/social-media-timeline/services/user/pkg/postgres"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	pgConfig := postgres.Config{
		Host:     config.PostgreSQL.Host,
		Port:     config.PostgreSQL.Port,
		User:     config.PostgreSQL.User,
		Password: config.PostgreSQL.Password,
		DBName:   config.PostgreSQL.DBName,
	}

	if err := postgres.InitPostgresClient(pgConfig); err != nil {
		log.Fatal("Failed to initialize Postgres client:", err)
	}

	pgClient := postgres.GetPostgresClient()

	/* mongo client initilization
	if err := mongodb.InitMongoClient(config.MongoDB.ConnectionString); err != nil {
		log.Fatal("Failed to initialize MongoDB client:", err)
	}

	mongoClient, err := mongodb.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	*/

	userRepository, err := repository.NewUserRepository("postgres", pgClient)
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewUserHandler(userRepository)
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
