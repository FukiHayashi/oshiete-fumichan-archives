package controllers

import (
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"takanome/database"
	"takanome/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 全tweetを表示
func TweetsHandler(ctx *gin.Context) {
	var tweets []models.Tweet

	// DB接続
	db := database.DataBaseConnect()
	defer database.DataBaseDisconnect(db)

	// ページ情報取得
	page := getPageInfo(ctx, &tweets, db)

	// tweetを取得
	db.Preload("Tags").Order("id DESC").Scopes(models.Paginate(page)).Find(&tweets)

	// 結果を返す
	ctx.HTML(http.StatusOK, "tweets.html", gin.H{
		"tweets": tweets,
		"page":   page,
	})
}

// ワードを含むtweetを検索
func TweetsSearchHandler(ctx *gin.Context) {
	var tweets []models.Tweet

	// ワードを受け取る
	search_words := strings.Split(ctx.PostForm("search_words"), " ")
	// DB接続
	db := database.DataBaseConnect()
	defer database.DataBaseDisconnect(db)

	// ページ情報取得
	page := getPageInfo(ctx, &tweets, db)

	// クエリ生成
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
	query.Scopes(models.Paginate(page)).Find(&tweets)

	// 結果を返す
	ctx.HTML(http.StatusOK, "tweets.html", gin.H{
		"tweets":           tweets,
		"search_words":     ctx.PostForm("search_words"),
		"search_condition": ctx.PostForm("search_condition"),
		"page":             page,
	})

}

// ページ情報取得
func getPageInfo(ctx *gin.Context, tweets *[]models.Tweet, db *gorm.DB) models.Page {
	total_elements := db.Find(tweets).RowsAffected

	page_number, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page_number == 0 || err != nil {
		page_number = 1
	}

	page_size, _ := strconv.Atoi(os.Getenv("PAGE_SIZE"))

	total_pages := int(math.Ceil(float64(total_elements) / float64(page_size)))
	paginate_link_num, _ := strconv.Atoi(os.Getenv("PAGINATE_LINK_NUM"))
	paginate_num := make([]int, paginate_link_num*2+1)

	for i := 0; i < len(paginate_num); i++ {
		paginate_num[i] = page_number - paginate_link_num + i
	}

	return models.Page{
		PageNumber:    page_number,
		PageSize:      page_size,
		TotalElements: total_elements,
		TotalPages:    total_pages,
		NextPage:      page_number + 1,
		PrevPage:      page_number - 1,
		PaginateNum:   paginate_num,
	}
}
