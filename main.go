package main

import (
	"log"
	"os"
	"takanome/database"
	"takanome/rareskill"
	"takanome/router"

	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .envファイルをロード
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// db初期化
	database.DataBaseInit()

	// ginをリリースモードに
	gin.SetMode(gin.ReleaseMode)
	rareskill.Register()

	jobrunner.Start()
	jobrunner.Schedule(os.Getenv("TWITTER_SEARCH_SCHEDULE"), rareskill.Skills{})

	rtr := router.New()
	rtr.Run(":3000")
}
