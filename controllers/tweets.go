package controllers

import (
	"net/http"
	"strings"
	"takanome/database"
	"takanome/models"

	"github.com/gin-gonic/gin"
)

// 全tweetを表示
func TweetsHandler(ctx *gin.Context) {
	// DB接続
	db := database.DataBaseConnect()
	defer database.DataBaseDisconnect(db)

	// tweetを取得
	var tweets []models.Tweet
	db.Preload("Tags").Order("id DESC").Find(&tweets)

	// 結果を返す
	ctx.HTML(http.StatusOK, "tweets.html", gin.H{
		"tweets": tweets,
	})
}

// ワードを含むtweetを検索
func TweetsSearchHandler(ctx *gin.Context) {
	// ワードを受け取る
	search_words := strings.Split(ctx.PostForm("search_words"), " ")
	// DB接続
	db := database.DataBaseConnect()
	defer database.DataBaseDisconnect(db)

	// クエリ生成
	var tweets []models.Tweet
	var query = db.Preload("Tags").Order("id DESC")
	for _, search_word := range search_words {
		w := "%" + search_word + "%"
		if ctx.PostForm("search_condition") == "AND" {
			query.Where("text LIKE ? OR retweet_text LIKE ?", w, w)
		} else {
			query.Or("text LIKE ? OR retweet_text LIKE ?", w, w)
		}
	}
	// tweetを取得
	query.Find(&tweets)

	// 結果を返す
	ctx.HTML(http.StatusOK, "tweets.html", gin.H{
		"tweets":           tweets,
		"search_words":     ctx.PostForm("search_words"),
		"search_condition": ctx.PostForm("search_condition"),
	})

}
