package main

import (
	"log"
	"takanome/database"
	"takanome/rareskill"
	"takanome/router"

	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
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
}

func main() {
	jobrunner.Start()
	jobrunner.Schedule("@every 1h", rareskill.Jobs{})

	rtr := router.New()
	rtr.Run(":3000")
}
