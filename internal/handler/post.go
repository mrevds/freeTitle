package handler

import (
	"github.com/gin-gonic/gin"
	"microblog/internal/model"
	"microblog/internal/repository"
	"net/http"
	"strconv"
)

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=255"`
	Content string `json:"content" binding:"required,min=1"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=255"`
	Content string `json:"content" binding:"required,min=1"`
}

func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// Получаем username из контекста
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Находим пользователя по username
	user, err := repository.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	// Создаем пост
	post := &model.Post{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: user.ID,
	}

	createdPost, err := repository.CreatePost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create post",
		})
		return
	}

	// Очищаем пароль автора в ответе
	createdPost.Author.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post":    createdPost,
	})
}

func GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	post, err := repository.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	// Очищаем пароль автора
	post.Author.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func GetAllPosts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	posts, err := repository.GetAllPosts(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch posts",
		})
		return
	}

	// Очищаем пароли авторов
	for i := range posts {
		posts[i].Author.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func GetMyPosts(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user, err := repository.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	posts, err := repository.GetPostsByAuthor(user.ID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch posts",
		})
		return
	}

	// Очищаем пароли авторов
	for i := range posts {
		posts[i].Author.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// Получаем username из контекста
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Находим пользователя
	user, err := repository.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	// Проверяем, что пост существует и принадлежит пользователю
	post, err := repository.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	if post.AuthorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only update your own posts",
		})
		return
	}

	// Обновляем пост
	updatedPost := &model.Post{
		Title:   req.Title,
		Content: req.Content,
	}

	result, err := repository.UpdatePost(id, updatedPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update post",
		})
		return
	}

	// Очищаем пароль автора
	result.Author.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"post":    result,
	})
}

func DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	// Получаем username из контекста
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Находим пользователя
	user, err := repository.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	// Проверяем, что пост существует и принадлежит пользователю
	post, err := repository.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	if post.AuthorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only delete your own posts",
		})
		return
	}

	// Удаляем пост
	if err := repository.DeletePost(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}
