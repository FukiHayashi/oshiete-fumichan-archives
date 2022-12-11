package main

import (
	"log"
	"takanome/database"
	"takanome/models"
	"takanome/rareskill"

	"github.com/joho/godotenv"
)

func init() {
	// .envファイルをロード
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// db初期化
	db := database.DataBaseConnect()
	db.AutoMigrate(&models.Tweet{})
	database.DataBaseDisconnect(db)
}

func main() {
	rareskill.Takanome()
}
