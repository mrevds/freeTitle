package repository

import (
	"microblog/internal/database"
	"microblog/internal/model"
)

func CreatePost(post *model.Post) (*model.Post, error) {
	result := database.DB.Create(post)
	if result.Error != nil {
		return nil, result.Error
	}

	database.DB.Preload("Author").First(post, post.ID)

	return post, nil
}

func GetPostByID(id int64) (*model.Post, error) {
	var post model.Post
	result := database.DB.Preload("Author").First(&post, id)
	if result.Error != nil {
		return nil, result.Error
	}

	count, _ := GetCommentsCountByPostID(post.ID)
	post.CommentsCount = count

	return &post, nil
}

func GetPostByIDWithComments(id int64, commentLimit, commentOffset int) (*model.Post, error) {
	var post model.Post
	result := database.DB.Preload("Author").First(&post, id)
	if result.Error != nil {
		return nil, result.Error
	}

	comments, err := GetCommentsByPostID(post.ID, commentLimit, commentOffset)
	if err != nil {
		return nil, err
	}

	for i := range comments {
		comments[i].Author.Password = ""
	}

	post.Comments = comments
	post.CommentsCount = int64(len(comments))

	return &post, nil
}

func GetAllPosts(limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	result := database.DB.Preload("Author").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	for i := range posts {
		count, _ := GetCommentsCountByPostID(posts[i].ID)
		posts[i].CommentsCount = count
	}

	return posts, nil
}

func GetPostsByAuthor(authorID int64, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	result := database.DB.Preload("Author").
		Where("author_id = ?", authorID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	for i := range posts {
		count, _ := GetCommentsCountByPostID(posts[i].ID)
		posts[i].CommentsCount = count
	}

	return posts, nil
}

func UpdatePost(id int64, post *model.Post) (*model.Post, error) {
	result := database.DB.Model(&model.Post{}).Where("id = ?", id).Updates(post)
	if result.Error != nil {
		return nil, result.Error
	}

	var updatedPost model.Post
	database.DB.Preload("Author").First(&updatedPost, id)

	return &updatedPost, nil
}

func DeletePost(id int64) error {
	result := database.DB.Delete(&model.Post{}, id)
	return result.Error
}
