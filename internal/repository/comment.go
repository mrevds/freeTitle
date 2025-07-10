package repository

import (
	"microblog/internal/database"
	"microblog/internal/model"
)

func CreateComment(comment *model.Comment) (*model.Comment, error) {
	result := database.DB.Create(comment)
	if result.Error != nil {
		return nil, result.Error
	}

	// Загружаем комментарий с автором
	database.DB.Preload("Author").First(comment, comment.ID)

	return comment, nil
}

func GetCommentsByPostID(postID int64, limit, offset int) ([]model.Comment, error) {
	var comments []model.Comment
	result := database.DB.Preload("Author").
		Where("post_id = ?", postID).
		Order("created_at asc").
		Limit(limit).
		Offset(offset).
		Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func GetCommentByID(id int64) (*model.Comment, error) {
	var comment model.Comment
	result := database.DB.Preload("Author").First(&comment, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

func UpdateComment(id int64, comment *model.Comment) (*model.Comment, error) {
	result := database.DB.Model(&model.Comment{}).Where("id = ?", id).Updates(comment)
	if result.Error != nil {
		return nil, result.Error
	}

	// Загружаем обновленный комментарий
	var updatedComment model.Comment
	database.DB.Preload("Author").First(&updatedComment, id)

	return &updatedComment, nil
}

func DeleteComment(id int64) error {
	result := database.DB.Delete(&model.Comment{}, id)
	return result.Error
}

func GetCommentsCountByPostID(postID int64) (int64, error) {
	var count int64
	result := database.DB.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&count)
	return count, result.Error
}
