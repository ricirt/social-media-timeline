package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ricirt/social-media-timeline/services/post/internal/model"
	"github.com/ricirt/social-media-timeline/services/post/internal/repository"
)

type PostHandler struct {
	repo repository.PostRepository
}

func NewPostHandler(repo repository.PostRepository) *PostHandler {
	return &PostHandler{
		repo: repo,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var post model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := post.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreatePost(c.Request.Context(), &post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	post, err := h.repo.GetPostByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	userID := c.Query("user_id")
	var (
		posts []model.Post
		err   error
	)

	if userID != "" {
		posts, err = h.repo.GetPostsByUserID(c.Request.Context(), userID)
	} else {
		posts, err = h.repo.GetPosts(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := post.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdatePost(c.Request.Context(), id, &post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeletePost(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *PostHandler) GetPostsByUserID(c *gin.Context) {
	userID := c.Param("userId")
	posts, err := h.repo.GetPostsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}
