package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ricirt/social-media-timeline/services/pkg/config"
	"github.com/ricirt/social-media-timeline/services/pkg/postgres"
	"github.com/ricirt/social-media-timeline/services/post/internal/handler"
	"github.com/ricirt/social-media-timeline/services/post/internal/repository"
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

	err = postgres.InitPostgresClient(pgConfig)
	if err != nil {
		log.Fatal("Failed to initialize Postgres client:", err)
	}

	pgClient := postgres.GetPostgresClient()

	postRepository, err := repository.NewPostRepository("postgres", pgClient)
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewPostHandler(postRepository)
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		posts := v1.Group("/posts")
		{
			posts.POST("", h.CreatePost)
			posts.GET("", h.GetPosts)
			posts.GET("/:id", h.GetPostByID)
			posts.PUT("/:id", h.UpdatePost)
			posts.DELETE("/:id", h.DeletePost)
			posts.GET("/users/:userId", h.GetPostsByUserID)
		}
	}

	r.Run(":" + config.Server.Port)
}
