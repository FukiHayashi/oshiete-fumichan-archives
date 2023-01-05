package controllers

import (
	"net/http"
	"takanome/database"
	"takanome/models"

	"github.com/gin-gonic/gin"
)

func IndexHandler(ctx *gin.Context) {
	// DB接続
	db := database.DataBaseConnect()
	defer database.DataBaseDisconnect(db)

	var tweet models.Tweet

	db.Order("id DESC").First(&tweet)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"categories": GetAllCategories(db),
		"tweeted_at": tweet.TweetedAt.Local().Format("2006/01/02 15:04:05"),
	})

}
