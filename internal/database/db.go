package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"microblog/internal/config"
	"microblog/internal/model"
	"time"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	dsn := cfg.GetDatabaseDSN()

	var db *gorm.DB
	var err error

	maxAttempts := 10
	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Database connection failed (attempt %d/%d): %v", i, maxAttempts, err)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		log.Fatalf("Could not connect to the database after %d attempts: %v", maxAttempts, err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	if err != nil {
		log.Fatal("Error in migration: ", err)
	}

	DB = db
}
