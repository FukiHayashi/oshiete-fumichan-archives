package controllers

import (
	"net/http"
	"takanome/database"
	"takanome/models"

	"github.com/gin-gonic/gin"
)

func TweetsHandler(ctx *gin.Context) {
	db := database.DataBaseConnect()
	defer database.DataBaseDisconnect(db)

	var tweets []models.Tweet
	db.Preload("Tags").Order("id DESC").Find(&tweets)
	ctx.HTML(http.StatusOK, "tweets.html", gin.H{
		"tweets": tweets,
	})
}
