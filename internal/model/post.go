package model

import (
	"time"
)

type Post struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Title         string    `json:"title" gorm:"size:255;not null"`
	Content       string    `json:"content" gorm:"type:text;not null"`
	AuthorID      int64     `json:"author_id" gorm:"not null"`
	Author        User      `json:"author" gorm:"foreignKey:AuthorID"`
	Comments      []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	CommentsCount int64     `json:"comments_count" gorm:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
