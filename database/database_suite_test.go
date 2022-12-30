package database_test

import (
	"log"
	"takanome/database"
	"takanome/models"
	"testing"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestDatabase(t *testing.T) {
	// テスト用設定ファイル読み込み
	err := godotenv.Load(".testenv")
	if err != nil {
		log.Fatalf("設定ファイルが見つかりません: %s", err)
	}
	// ginkgo
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database Suite")
}

var _ = Describe("database", Ordered, func() {
	Describe("func DataBaseConnect", func() {
		Context(".envのDATABASE_URLに接続できた時", func() {
			var db *gorm.DB
			AfterAll(func() {
				database.DataBaseDisconnect(db)
			})
			It("コネクションが返ること", func() {
				db = database.DataBaseConnect()
				Expect(db).ToNot(BeNil())
			})
		})
	})
	Describe("func DataBaseInit", Ordered, func() {
		var db *gorm.DB
		BeforeAll(func() {
			db = database.DataBaseConnect()
			// 開始前にDBのテーブルを全てドロップする
			db.Migrator().DropTable(&models.Tweet{}, &models.Category{}, &models.Group{}, &models.Tag{}, &models.Keyword{}, &models.History{}, &models.TweetTags{})
		})
		AfterAll(func() {
			// 終了後にDBのテーブルを全てドロップする
			db.Migrator().DropTable(&models.Tweet{}, &models.Category{}, &models.Group{}, &models.Tag{}, &models.Keyword{}, &models.History{}, &models.TweetTags{})

			database.DataBaseDisconnect(db)
		})
		Context("実行履歴がない時", func() {
			It("DBがCSVの内容で初期化されること", func() {
				database.DataBaseInit()

				var category models.Category
				db.First(&category)
				Expect(category.Name).To(Equal("Character"))

				var group models.Group
				db.First(&group)
				Expect(group.Name).To(Equal("百合ヶ丘女学院"))

				var tag models.Tag
				db.First(&tag)
				Expect(tag.Name).To(Equal("青木夏帆"))

				var keyword models.Keyword
				db.First(&keyword)
				Expect(keyword.Name).To(Equal("青木夏帆"))

				var registered_at models.History
				db.First(&registered_at, "name = ?", "registeredAt")
				Expect(registered_at.Name).To(Equal("registeredAt"))
			})
		})
		Context("実行履歴がある時", func() {
			It("DBが初期化されないこと", func() {
				var before_registered_at models.History
				db.First(&before_registered_at, "name = ?", "registeredAt")

				database.DataBaseInit()

				var after_registered_at models.History
				db.First(&after_registered_at, "name=?", "registeredAt")

				// registeredAtのLastUpdatedAtが同じであればDBは初期化されていない
				Expect(before_registered_at.LastUpdatedAt).To(Equal(after_registered_at.LastUpdatedAt))

			})
		})
	})
})
