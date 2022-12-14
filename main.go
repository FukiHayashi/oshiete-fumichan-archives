package main

import (
	"log"
	"net/http"
	"takanome/database"
	"takanome/rareskill"

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
}

func main() {
	jobrunner.Start()
	jobrunner.Schedule("@every 2h", rareskill.Jobs{})

	router := gin.Default()
	router.GET("/jobrunner/status", JobResult)
	router.Run(":3000")
}

func JobResult(c *gin.Context) {
	c.JSON(http.StatusOK, jobrunner.StatusJson())
}
