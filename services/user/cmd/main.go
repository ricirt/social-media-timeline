package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ricirt/social-media-timeline/services/pkg/config"
	"github.com/ricirt/social-media-timeline/services/pkg/postgres"
	"github.com/ricirt/social-media-timeline/services/user/internal/handler"
	"github.com/ricirt/social-media-timeline/services/user/internal/repository"

	//mongodb "github.com/ricirt/social-media-timeline/services/user/pkg/mongo-db"
	"log"
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

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", h.CreateUser)
			users.GET("", h.GetUsers)
			users.GET("/:id", h.GetUserByID)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
		}
	}

	r.Run(":" + config.Server.Port)
}
