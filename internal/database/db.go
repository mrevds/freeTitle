package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"microblog/internal/model"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=denis dbname=microblog password=11042005 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error in connection", err)
	}
	err = db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	if err != nil {
		log.Fatal("Error in migration", err)
	}
	DB = db
}
