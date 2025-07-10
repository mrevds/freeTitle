package model

import (
	"time"
)

type Comment struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	PostID    int64     `json:"post_id" gorm:"not null"`
	Post      Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	AuthorID  int64     `json:"author_id" gorm:"not null"`
	Author    User      `json:"author" gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
