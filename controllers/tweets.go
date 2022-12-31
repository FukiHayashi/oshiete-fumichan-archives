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

	// 検索キーワードを受け取る
	keywords := ctx.DefaultQuery("keywords", "")
	search_words := strings.Split(keywords, " ")

	// クエリ生成
	var query = db.Preload("Tags").Order("id DESC")
	for _, search_word := range search_words {
		w := "%" + search_word + "%"
		query.Or("text LIKE ? OR retweet_text LIKE ?", w, w)
	}
	// ページ情報取得
	page := getPageInfo(ctx, &tweets, query)

	// tweetを取得
	query.Scopes(models.Paginate(page)).Find(&tweets)

	// 結果を返す
	ctx.HTML(http.StatusOK, "tweets.html", gin.H{
		"tweets":   tweets,
		"keywords": keywords,
		"page":     page,
	})
}

// ページ情報取得
func getPageInfo(ctx *gin.Context, tweets *[]models.Tweet, db *gorm.DB) models.Page {
	// 合計の要素数
	total_elements := db.Find(tweets).RowsAffected

	// 現在ページ
	page_number, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page_number == 0 || err != nil {
		page_number = 1
	}

	// １ページあたりの要素数
	page_size, _ := strconv.Atoi(os.Getenv("PAGE_SIZE"))

	// 合計ページ
	total_pages := int(math.Ceil(float64(total_elements) / float64(page_size)))

	// 画面に表示するページリンクの数
	paginate_link_num, _ := strconv.Atoi(os.Getenv("PAGINATE_LINK_NUM"))

	// ページネーションの情報
	var pageinate_infos []models.PaginateInfo

	if total_elements > 0 {
		// キーワードがある場合に情報を追加する
		keywords := ctx.DefaultQuery("keywords", "")
		var keywords_path_param = ""

		if keywords != "" {
			keywords_path_param = "&keywords=" + strings.ReplaceAll(keywords, " ", "+")
		}

		// 最初のページ情報
		pageinate_infos = append(pageinate_infos, models.PaginateInfo{
			PageNumber: 1, PathParam: "?page=1" + keywords_path_param, Info: "first"})

		// 画面に表示するページの情報
		for i := paginate_link_num * (-1); i <= paginate_link_num; i++ {
			number := page_number + i
			if 0 < number && number <= total_pages {
				pageinate_infos = append(pageinate_infos, models.PaginateInfo{
					PageNumber: number, PathParam: "?page=" + strconv.Itoa(number) + keywords_path_param, Info: ""})
			}
		}

		// 最後のページ情報
		pageinate_infos = append(pageinate_infos, models.PaginateInfo{
			PageNumber: total_pages, PathParam: "?page=" + strconv.Itoa(total_pages) + keywords_path_param, Info: "last"})
	}
	return models.Page{
		PageNumber:    page_number,
		PageSize:      page_size,
		TotalElements: total_elements,
		TotalPages:    total_pages,
		PaginateInfos: pageinate_infos,
	}
}
