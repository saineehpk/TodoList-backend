package database

import (
	"log"

	"backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err := DB.AutoMigrate(&models.User{}, &models.Board{}, &models.Card{}, &models.Label{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
