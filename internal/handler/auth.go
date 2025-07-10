package handler

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"microblog/internal/model"
	"microblog/internal/repository"
	"microblog/internal/util"
	"net/http"
	"regexp"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         *model.User `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,}$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	if len(req.Username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username must be at least 3 characters long",
		})
		return
	}

	if !usernameRegex.MatchString(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username must contain only Latin letters, numbers, underscores, and hyphens",
		})
		return
	}

	if !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email format",
		})
		return
	}

	if user, _ := repository.GetUserByUsername(req.Username); user != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username already taken",
		})
		return
	}

	if user, _ := repository.GetUserByEmail(req.Email); user != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already registered",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process password",
		})
		return
	}

	user, err := repository.CreateUser(&model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	accessToken, err := util.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	refreshToken, err := util.GenerateRefreshToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate refresh token",
		})
		return
	}

	// Сохраняем refresh token в базе данных
	if err := repository.UpdateUserRefreshToken(user.Username, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save refresh token",
		})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// Валидируем refresh token
	_, err := util.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid refresh token",
		})
		return
	}

	// Проверяем, что refresh token существует в базе данных
	user, err := repository.GetUserByRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token not found or expired",
		})
		return
	}

	// Генерируем новые токены
	newAccessToken, err := util.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate new access token",
		})
		return
	}

	newRefreshToken, err := util.GenerateRefreshToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate new refresh token",
		})
		return
	}

	// Обновляем refresh token в базе данных
	if err := repository.UpdateUserRefreshToken(user.Username, newRefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}
func Logout(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Очищаем refresh token из базы данных
	if err := repository.ClearUserRefreshToken(username.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
