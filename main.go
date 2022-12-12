package main

import (
	"log"
	"net/http"
	"takanome/database"
	"takanome/models"
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
	db := database.DataBaseConnect()
	db.AutoMigrate(&models.Tweet{}, &models.Category{}, &models.Tag{}, &models.Keyword{})
	database.DataBaseDisconnect(db)
}

func main() {
	jobrunner.Start()
	jobrunner.Schedule("@every 2h", rareskill.JobTakanome{})

	router := gin.Default()
	router.GET("/jobrunner/status", JobResult)
	router.Run(":3000")
}

func JobResult(c *gin.Context) {
	c.JSON(http.StatusOK, jobrunner.StatusJson())
}
