package repository

import (
	"microblog/internal/database"
	"microblog/internal/model"
	"time"
)

func CreateUser(user *model.User) (*model.User, error) {
	result := database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func UpdateUserRefreshToken(username, refreshToken string) error {
	result := database.DB.Model(&model.User{}).
		Where("username = ?", username).
		Updates(map[string]interface{}{
			"refresh_token": refreshToken,
			"token_expiry":  time.Now().Add(7 * 24 * time.Hour),
		})
	return result.Error
}

func GetUserByRefreshToken(refreshToken string) (*model.User, error) {
	var user model.User
	result := database.DB.Where("refresh_token = ? AND token_expiry > ?", refreshToken, time.Now()).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func ClearUserRefreshToken(username string) error {
	result := database.DB.Model(&model.User{}).
		Where("username = ?", username).
		Updates(map[string]interface{}{
			"refresh_token": "",
			"token_expiry":  time.Time{},
		})
	return result.Error
}
