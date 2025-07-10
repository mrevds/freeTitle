package model

import "time"

type User struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string    `json:"username" gorm:"size:100;not null;uniqueIndex"`
	Password     string    `json:"password" gorm:"size:255;not null"`
	Email        string    `json:"email" gorm:"size:100;not null;uniqueIndex"`
	RefreshToken string    `json:"-" gorm:"size:500"`
	TokenExpiry  time.Time `json:"-"`
}
