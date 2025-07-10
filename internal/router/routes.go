package router

import (
	"github.com/gin-gonic/gin"
	"microblog/internal/handler"
	"microblog/internal/middleware"
)

func Routers() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", handler.Ping)

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.RefreshToken)
		auth.POST("/logout", middleware.AuthMiddleware(), handler.Logout)
	}

	// Публичные маршруты для постов
	posts := r.Group("/api/posts")
	{
		posts.GET("", handler.GetAllPosts)                           // GET /api/posts
		posts.GET("/:id", handler.GetPost)                           // GET /api/posts/1
		posts.GET("/:id/with-comments", handler.GetPostWithComments) // GET /api/posts/1/with-comments
		posts.GET("/:id/comments", handler.GetCommentsByPost)        // GET /api/posts/1/comments
	}

	// Защищенные маршруты
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Маршруты для постов (требуют авторизации)
		api.POST("/posts", handler.CreatePost)       // POST /api/posts
		api.GET("/posts/my", handler.GetMyPosts)     // GET /api/posts/my
		api.PUT("/posts/:id", handler.UpdatePost)    // PUT /api/posts/1
		api.DELETE("/posts/:id", handler.DeletePost) // DELETE /api/posts/1

		// Маршруты для комментариев (требуют авторизации)
		api.POST("/posts/:id/comments", handler.CreateComment) // POST /api/posts/1/comments
		api.PUT("/comments/:id", handler.UpdateComment)        // PUT /api/comments/1
		api.DELETE("/comments/:id", handler.DeleteComment)     // DELETE /api/comments/1
	}

	return r
}
