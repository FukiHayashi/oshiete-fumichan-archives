package main_test

import (
	"log"
	"net/http/httptest"
	"os"
	"takanome/database"
	"takanome/models"
	"takanome/router"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	"gorm.io/gorm"
)

func TestTakanome(t *testing.T) {
	// ginkgo
	RegisterFailHandler(Fail)
	RunSpecs(t, "Takanome Suite")
}

var (
	agouti_driver *agouti.WebDriver
	test_server   *httptest.Server
	db            *gorm.DB
	server_url    string
)

var _ = BeforeSuite(func() {
	// agoutiドライバの初期化
	agouti_driver = agouti.ChromeDriver()
	Expect(agouti_driver.Start()).To(Succeed())

	println(os.Getenv("DATABASE_URL"))

	// テスト用設定ファイル読み込み
	err := godotenv.Load(".testenv")
	if err != nil {
		log.Fatalf("設定ファイルが見つかりません: %s", err)
	}

	println(os.Getenv("DATABASE_URL"))

	// DBコネクション
	db = database.DataBaseConnect()

	// DB初期化
	db.Migrator().DropTable(&models.Tweet{}, &models.Category{}, &models.Group{}, &models.Tag{}, &models.Keyword{}, &models.History{}, &models.TweetTags{})
	database.DataBaseInit()

	// テストサーバを立てる
	gin.SetMode(gin.TestMode)
	test_server = httptest.NewServer(router.New())
	server_url = test_server.URL
})

var _ = AfterSuite(func() {
	// agoutiドライバ停止
	Expect(agouti_driver.Stop()).To(Succeed())

	// 終了後にDBのテーブルを全てドロップする
	db.Migrator().DropTable(&models.Tweet{}, &models.Category{}, &models.Group{}, &models.Tag{}, &models.Keyword{}, &models.History{}, &models.TweetTags{})

	// DBコネクション削除
	database.DataBaseDisconnect(db)

	// テストサーバ停止
	test_server.Close()
})
