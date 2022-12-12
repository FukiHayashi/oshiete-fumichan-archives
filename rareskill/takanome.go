package rareskill

import (
	"encoding/json"
	"log"
	"os"
	"takanome/database"
	"takanome/models"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type JobTakanome struct{}

func (e JobTakanome) Run() {
	Takanome()
	log.Println(time.Now(), "鷹の目実行")
}

// 二水ちゃんのツイートをDBへ保存
func Takanome() {
	// アカウントのファボのツイートを取得
	tlist := getFavoriteList()
	// DBへ格納するためのモデルへ変換
	mtl := tweetList2modelList(tlist)
	// DBへ格納
	writeToDB(mtl)
}

// DBへデータを書き込み
func writeToDB(mtl []*models.Tweet) {
	db := database.DataBaseConnect()
	err := db.Save(&mtl)
	if err.Error != nil {
		log.Println(err.Error.Error())
	}
	database.DataBaseDisconnect(db)
}

// ファボのツイートをリストで取得
func getFavoriteList() []twitter.Tweet {
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_SECRET_KEY"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_TOKEN_SECRET"))

	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// twitter client
	client := twitter.NewClient(httpClient)

	// TWITTER_ACCOUNTのファボを全文で取得する設定
	params := twitter.FavoriteListParams{
		ScreenName: os.Getenv("TWITTER_ACCOUNT"),
		TweetMode:  "extended", // tweet全文で取得
	}

	// ファボのリストを取得
	tlist, _, _ := client.Favorites.List(&params)
	return tlist
}

// tweetの配列をtweetモデル配列に変換
func tweetList2modelList(twl []twitter.Tweet) (mtl []*models.Tweet) {
	for _, val := range twl {
		mt := tweet2model(&val)
		mtl = append(mtl, mt)
	}
	return mtl
}

// tweetをtweetモデルに変換
func tweet2model(tw *twitter.Tweet) *models.Tweet {
	// ツイートの情報
	mt := models.Tweet{
		ID:         tw.ID,
		Text:       tw.FullText,
		TweetedAt:  str2time(tw.CreatedAt),
		Url:        "https://twitter.com/" + tw.User.ScreenName + "/status/" + tw.IDStr,
		RawData:    tweet2json(tw),
		ScreenName: tw.User.ScreenName,
	}
	// リツイートの情報
	if tw.QuotedStatus != nil {
		mt.RetweetScreenName = tw.QuotedStatus.User.ScreenName
		mt.RetweetText = tw.QuotedStatus.FullText
		mt.RetweetUrl = "https://twitter.com/" + tw.QuotedStatus.User.ScreenName + "/status/" + tw.QuotedStatus.IDStr
	}
	return &mt
}

// tweetをjsonへ変換
func tweet2json(tw *twitter.Tweet) string {
	jsonData, _ := json.Marshal(tw)
	return string(jsonData)
}

// tweetで取得した時間をtime型に変換
func str2time(t string) time.Time {
	parsedTime, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", t)
	return parsedTime
}
