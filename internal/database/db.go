package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"microblog/internal/config"
	"microblog/internal/model"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	dsn := cfg.GetDatabaseDSN()
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
