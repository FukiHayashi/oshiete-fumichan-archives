package rareskill_test

import (
	"log"
	"takanome/database"
	"takanome/models"
	"takanome/rareskill"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestRareskill(t *testing.T) {
	// テスト用設定ファイル読み込み
	err := godotenv.Load(".testenv")
	if err != nil {
		log.Fatalf("設定ファイルが見つかりません: %s", err)
	}
	// ginkgo
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rareskill Suite")
}

var _ = Describe("rareskill", Ordered, func() {
	var db *gorm.DB

	BeforeAll(func() {
		database.DataBaseInit()
		db = database.DataBaseConnect()

	})
	AfterAll(func() {
		// 終了後にDBのテーブルを全てドロップする
		db.Migrator().DropTable(&models.Tweet{}, &models.Category{}, &models.Group{}, &models.Tag{}, &models.Keyword{}, &models.History{}, &models.TweetTags{})

		database.DataBaseDisconnect(db)
	})

	Describe("func Takanome", func() {
		var before_tweets_count int
		var tweets []models.Tweet

		BeforeEach(func() {
			// 実行前のDB内のツイート数取得
			db.Find(&tweets)
			before_tweets_count = len(tweets)

			rareskill.Takanome()
		})

		Context("ツイートを取得できた時", func() {
			It("DBにツイートが登録されること", func() {
				db.Find(&tweets)
				Expect(len(tweets)).To(BeNumerically(">", before_tweets_count))
			})
		})
		Context("ツイートを取得できない時", func() {
			It("DBにツイートが登録されないこと", func() {
				db.Find(&tweets)
				Expect(len(tweets)).To(BeNumerically("==", before_tweets_count))
			})
		})
	})
	Describe("func Register", Ordered, func() {
		var include_keyword_tweet, not_include_keyword_tweet models.Tweet
		BeforeAll(func() {
			include_keyword_tweet.ID = int64(1)
			include_keyword_tweet.Text = "梨璃"
			include_keyword_tweet.TweetedAt = time.Now()
			include_keyword_tweet.Url = "http://localhost"
			include_keyword_tweet.RawData = "include"
			include_keyword_tweet.ScreenName = "include"

			not_include_keyword_tweet.ID = int64(2)
			not_include_keyword_tweet.Text = "not include"
			not_include_keyword_tweet.TweetedAt = time.Now()
			not_include_keyword_tweet.Url = "http://localhost"
			not_include_keyword_tweet.RawData = "not_include"
			not_include_keyword_tweet.ScreenName = "not_include"

			db.Save(&include_keyword_tweet)
			db.Save(&not_include_keyword_tweet)
		})

		Context("最終実行時刻より新しいツイートがDBにある時", func() {
			var before_last_updated_at models.History
			BeforeAll(func() {
				db.First(&before_last_updated_at, "name = ?", "registeredAt")
				rareskill.Register()
			})
			It("キーワードを含むツイートにタグが設定されること", func() {
				var tweet models.Tweet
				db.Preload("Tags").First(&tweet, include_keyword_tweet.ID)
				Expect(len(tweet.Tags)).To(Equal(1))
			})
			It("キーワードを含まないツイートにタグが設定されないこと", func() {
				var tweet models.Tweet
				db.Preload("Tags").First(&tweet, not_include_keyword_tweet.ID)
				Expect(len(tweet.Tags)).To(Equal(0))
			})
			It("最終実行時刻が更新されること", func() {
				var last_updated_at models.History
				db.First(&last_updated_at, "name = ?", "registeredAt")
				Expect(last_updated_at.LastUpdatedAt).ToNot(Equal(before_last_updated_at.LastUpdatedAt))
			})
		})
		Context("最終実行時刻より新しいツイートがDBにない時", func() {
			var before_last_updated_at models.History
			BeforeAll(func() {
				db.First(&before_last_updated_at, "name = ?", "registeredAt")
				rareskill.Register()
			})

			It("キーワードを含むツイートが更新されないこと", func() {
				var tweet models.Tweet
				db.Preload("Tags").First(&tweet, include_keyword_tweet.ID)
				Expect(len(tweet.Tags)).To(Equal(1))
			})
			It("キーワードを含まないツイートが更新されないこと", func() {
				var tweet models.Tweet
				db.Preload("Tags").First(&tweet, not_include_keyword_tweet.ID)
				Expect(len(tweet.Tags)).To(Equal(0))
			})
			It("最終実行時刻が更新されないこと", func() {
				var last_updated_at models.History
				db.First(&last_updated_at, "name = ?", "registeredAt")
				Expect(last_updated_at.LastUpdatedAt).ToNot(Equal(before_last_updated_at.LastUpdatedAt))
			})
		})
	})
})
