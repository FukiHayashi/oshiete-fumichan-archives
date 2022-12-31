package main_test

import (
	"os"
	"strconv"
	"takanome/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("/tweets", Ordered, func() {
	var (
		page             *agouti.Page
		before_page_size string
		tweets           []models.Tweet
	)

	BeforeAll(func() {
		before_page_size = os.Getenv("PAGE_SIZE")
		// １ページに表示するツイート数を強制的に1にする
		os.Setenv("PAGE_SIZE", "1")
		for i := 0; i < 10; i++ {
			tweet := tweetFactory.MustCreate().(*models.Tweet)
			tweets = append(tweets, *tweet)
		}
		db.Save(&tweets)
	})
	AfterAll(func() {
		db.Delete(&tweets)
		// 設定を元に戻す
		os.Setenv("PAGE_SIZE", before_page_size)
	})
	BeforeEach(func() {
		var err error
		page, err = agouti_driver.NewPage()
		Expect(err).ToNot(HaveOccurred())
	})
	AfterEach(func() {
		Expect(page.Destroy()).To(Succeed())
	})
	Describe("ページコンテンツ確認", func() {
		BeforeEach(func() {
			Expect(page.Navigate(server_url + "/tweets")).To(Succeed())
		})
		Describe("表示内容確認", func() {
			Context("ページを開いた時", func() {
				It("tweetsページが表示されること", func() {
					Expect(page).To(HaveTitle("教えて二水ちゃん | tweets"))
				})
				It("最新のツイートが表示されること", func() {
					var tweet models.Tweet
					db.Order("id DESC").First(&tweet)
					Expect(page.Find("#tweet-" + strconv.FormatInt(tweet.ID, 10))).To(HaveText(tweet.Text))
				})
				It("表示されるツイートの数が.envに設定した値と等しいこと", func() {
					num, _ := strconv.Atoi(os.Getenv("PAGE_SIZE"))
					Expect(page.Find("#tweet").Count()).To(Equal(num))
				})
			})
		})
		Describe("ページング確認", func() {
			Context("次ページボタンを押した時", func() {
				BeforeEach(func() {
					Expect(page.FindByID("btn-page-2").Click()).To(Succeed())
				})
				It("次のページが表示されること", func() {
					Expect(page).To(HaveURL(server_url + "/tweets?page=2"))
				})
				It("最新から２番目のツイートが表示されること", func() {
					// 配列の最後が一番新しいツイートなので、後ろから２番目が該当のツイート
					tweet := tweets[len(tweets)-2]
					Expect(page.Find("#tweet-" + strconv.FormatInt(tweet.ID, 10))).To(HaveText(tweet.Text))
				})
			})
			Context("最後ボタンを押した時", func() {
				BeforeEach(func() {
					Expect(page.FindByID("btn-page-last").Click()).To(Succeed())
				})
				It("最後のページが表示されること", func() {
					Expect(page).To(HaveURL(server_url + "/tweets?page=" + strconv.Itoa(len(tweets))))
				})
				It("最古のツイートが表示されること", func() {
					tweet := tweets[0]
					Expect(page.Find("#tweet-" + strconv.FormatInt(tweet.ID, 10))).To(HaveText(tweet.Text))
				})
			})
			Context("前ページボタンを押した時", func() {
				BeforeEach(func() {
					Expect(page.Navigate(server_url + "/tweets?page=" + strconv.Itoa(len(tweets)))).To(Succeed())
					Expect(page.FindByID("btn-page-" + strconv.Itoa(len(tweets)-1)).Click()).To(Succeed())
				})
				It("前のページが表示されること", func() {
					Expect(page).To(HaveURL(server_url + "/tweets?page=" + strconv.Itoa(len(tweets)-1)))
				})
				It("最古から２番目のツイートが表示されること", func() {
					tweet := tweets[1]
					Expect(page.Find("#tweet-" + strconv.FormatInt(tweet.ID, 10))).To(HaveText(tweet.Text))
				})
			})
			Context("最初ボタンを押した時", func() {
				BeforeEach(func() {
					Expect(page.Navigate(server_url + "/tweets?page=" + strconv.Itoa(len(tweets)))).To(Succeed())
					Expect(page.FindByID("btn-page-first").Click()).To(Succeed())
				})
				It("最初のページが表示されること", func() {
					Expect(page).To(HaveURL(server_url + "/tweets?page=1"))
				})
				It("最新のツイートが表示されること", func() {
					tweet := tweets[len(tweets)-1]
					Expect(page.Find("#tweet-" + strconv.FormatInt(tweet.ID, 10))).To(HaveText(tweet.Text))
				})
			})
		})
		Describe("検索機能確認", func() {
			Context("該当するツイートがあるキーワードで検索した時", func() {
				BeforeEach(func() {
					Expect(page.FindByID("input-search-keywords").Fill("tweet-1 tweet-10")).To(Succeed())
					Expect(page.FindByID("btn-search").Click()).To(Succeed())
				})
				It("該当ツイートが表示されること", func() {
					Expect(page.FindByID("tweet-10")).To(HaveText("tweet-10"))
					Expect(page.FindByID("btn-page-2").Click()).To(Succeed())
					Expect(page.FindByID("tweet-1")).To(HaveText("tweet-1"))
				})
			})
			Context("該当するツイートがないキーワードで検索した時", func() {
				BeforeEach(func() {
					Expect(page.FindByID("input-search-keywords").Fill("no-word")).To(Succeed())
					Expect(page.FindByID("btn-search").Click()).To(Succeed())
				})
				It("ツイートが見つかりませんと表示されること", func() {
					Expect(page.FindByID("tweet-not-found")).To(HaveText("ツイートが見つかりません"))
				})
			})
		})
	})
	Describe("URLパラメータ確認", func() {
		Describe("?page=", func() {
			Context("1ページ目を開いた時", func() {
				BeforeEach(func() {
					Expect(page.Navigate(server_url + "/tweets?page=1")).To(Succeed())
				})
				It("最新のツイートが表示されること", func() {
					tweet := tweets[len(tweets)-1]
					Expect(page.Find("#tweet-" + strconv.FormatInt(tweet.ID, 10))).To(HaveText(tweet.Text))
				})
			})
			Context("最終ページを開いた時", func() {
				BeforeAll(func() {
					tweet_num := strconv.Itoa(len(tweets))
					Expect(page.Navigate(server_url + "/tweets?page=" + tweet_num)).To(Succeed())
				})
				It("最古のツイートが表示されること", func() {
					tweet := tweets[0]
					Expect(page.Find("#tweet-" + strconv.FormatInt(tweet.ID, 10))).To(HaveText(tweet.Text))
				})
			})
		})
		Describe("?keyword=", func() {
			Context("該当するツイートがあるキーワードで検索した時", func() {
				BeforeEach(func() {
					Expect(page.Navigate(server_url + "/tweets?keywords=tweet-10+tweet-1")).To(Succeed())
				})
				It("該当ツイートが表示されること", func() {
					Expect(page.FindByID("tweet-10")).To(HaveText("tweet-10"))
					Expect(page.FindByID("btn-page-2").Click()).To(Succeed())
					Expect(page.FindByID("tweet-1")).To(HaveText("tweet-1"))
				})
			})
			Context("該当するツイートがないキーワードで検索した時", func() {
				BeforeEach(func() {
					Expect(page.Navigate(server_url + "/tweets?keywords=no-word")).To(Succeed())
				})
				It("ツイートが見つかりませんと表示されること", func() {
					Expect(page.FindByID("tweet-not-found")).To(HaveText("ツイートが見つかりません"))
				})
			})
		})
	})
})
