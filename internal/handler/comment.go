package handler

import (
	"github.com/gin-gonic/gin"
	"microblog/internal/model"
	"microblog/internal/repository"
	"net/http"
	"strconv"
)

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

func CreateComment(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

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

	_, err = repository.GetPostByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	comment := &model.Comment{
		Content:  req.Content,
		PostID:   postID,
		AuthorID: user.ID,
	}

	createdComment, err := repository.CreateComment(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create comment",
		})
		return
	}

	createdComment.Author.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": createdComment,
	})
}

func GetCommentsByPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	_, err = repository.GetPostByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	comments, err := repository.GetCommentsByPostID(postID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch comments",
		})
		return
	}

	for i := range comments {
		comments[i].Author.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

func GetPostWithComments(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	commentLimitStr := c.DefaultQuery("comment_limit", "10")
	commentOffsetStr := c.DefaultQuery("comment_offset", "0")

	commentLimit, err := strconv.Atoi(commentLimitStr)
	if err != nil || commentLimit < 1 || commentLimit > 50 {
		commentLimit = 10
	}

	commentOffset, err := strconv.Atoi(commentOffsetStr)
	if err != nil || commentOffset < 0 {
		commentOffset = 0
	}

	post, err := repository.GetPostByIDWithComments(postID, commentLimit, commentOffset)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	post.Author.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func UpdateComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

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

	comment, err := repository.GetCommentByID(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Comment not found",
		})
		return
	}

	if comment.AuthorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only update your own comments",
		})
		return
	}

	updatedComment := &model.Comment{
		Content: req.Content,
	}

	result, err := repository.UpdateComment(commentID, updatedComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update comment",
		})
		return
	}

	result.Author.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
		"comment": result,
	})
}

func DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

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

	comment, err := repository.GetCommentByID(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Comment not found",
		})
		return
	}

	if comment.AuthorID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only delete your own comments",
		})
		return
	}

	if err := repository.DeleteComment(commentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete comment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}
