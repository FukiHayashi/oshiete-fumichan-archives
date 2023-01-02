package main_test

import (
	"net/url"
	"os"
	"strconv"
	"strings"
	"takanome/models"
	"takanome/rareskill"
	"time"

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
		// ツイートを作成
		for i := 0; i < 10; i++ {
			tweet := tweetFactory.MustCreate().(*models.Tweet)
			tweets = append(tweets, *tweet)
		}
		db.Save(&tweets)
		// タグを登録
		rareskill.Register()
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
					Expect(page.FindByID("search-menue").Click()).To(Succeed())
					time.Sleep(time.Second * 2) // 入力画面が出るまで待つ
					Expect(page.FindByID("input-search-keywords").Fill("tweet-1 tweet-10")).To(Succeed())
					Expect(page.FindByID("btn-search").Click()).To(Succeed())
				})
				It("該当ツイートが表示されること", func() {
					Expect(page.FindByID("tweet-10")).To(MatchText("tweet-10"))
					Expect(page.FindByID("btn-page-2").Click()).To(Succeed())
					Expect(page.FindByID("tweet-1")).To(MatchText("tweet-1"))
				})
			})
			Context("該当するツイートがないキーワードで検索した時", func() {
				BeforeEach(func() {
					Expect(page.FindByID("search-menue").Click()).To(Succeed())
					time.Sleep(time.Second * 2) // 入力画面が出るまで待つ
					Expect(page.FindByID("input-search-keywords").Fill("no-word")).To(Succeed())
					Expect(page.FindByID("btn-search").Click()).To(Succeed())
				})
				It("ツイートが見つかりませんと表示されること", func() {
					Expect(page.FindByID("tweet-not-found")).To(HaveText("ツイートが見つかりません"))
				})
			})
		})
		Describe("タグ検索機能", func() {
			Context("該当するツイートがあるタグで検索した時", func() {
				var (
					tag         models.Tag
					tag_id      string
					group_id    string
					category_id string
				)
				BeforeEach(func() {
					db.Preload("Group").Where("name = ?", "二川二水").Find(&tag)
					tag_id = strconv.FormatInt(int64(tag.ID), 10)
					group_id = strconv.FormatInt(int64(tag.GroupID), 10)
					category_id = strconv.FormatInt(int64(tag.Group.CategoryID), 10)

					Expect(page.FindByID("search-menue").Click()).To(Succeed())            //検索メニューを開く
					time.Sleep(time.Second * 1)                                            // 入力画面が出るまで待つ
					Expect(page.FindByID("category-" + category_id).Click()).To(Succeed()) // カテゴリーを開く
					time.Sleep(time.Second * 1)                                            // 入力画面が出るまで待つ

					Expect(page.FindByID("group-" + group_id).Click()).To(Succeed()) // グループを開く
					time.Sleep(time.Second * 1)                                      // 入力画面が出るまで待つ

					Expect(page.FindByID("tag-" + tag_id).Click()).To(Succeed()) // タグをクリック
				})
				It("該当ツイートが表示されること", func() {
					Expect(page).To(HaveURL(server_url + "/tweets/" + strings.ToLower(url.PathEscape(tag.Name))))

					var tweet_id = strconv.FormatInt(tweets[len(tweets)-1].ID, 10)
					Expect(page.FindByID("tweet-" + tweet_id)).To(MatchText(tag.Name))
					Expect(page.FindByID("tag-" + tweet_id + "-" + tag_id)).To(MatchText(tag.Name))
					Expect(page.FindByID("btn-page-2").Click()).To(Succeed())

					tweet_id = strconv.FormatInt(tweets[len(tweets)-2].ID, 10)
					Expect(page.FindByID("tweet-" + tweet_id)).To(MatchText(tag.Name))
					Expect(page.FindByID("tag-" + tweet_id + "-" + tag_id)).To(MatchText(tag.Name))
				})
			})
			Context("該当するツイートがないタグで検索した時", func() {
				var (
					tag         models.Tag
					tag_id      string
					group_id    string
					category_id string
				)
				BeforeEach(func() {
					db.Preload("Group").Where("name = ?", "一柳梨璃").Find(&tag)
					tag_id = strconv.FormatInt(int64(tag.ID), 10)
					group_id = strconv.FormatInt(int64(tag.GroupID), 10)
					category_id = strconv.FormatInt(int64(tag.Group.CategoryID), 10)

					Expect(page.FindByID("search-menue").Click()).To(Succeed())            //検索メニューを開く
					time.Sleep(time.Second * 1)                                            // 入力画面が出るまで待つ
					Expect(page.FindByID("category-" + category_id).Click()).To(Succeed()) // カテゴリーを開く
					time.Sleep(time.Second * 1)                                            // 入力画面が出るまで待つ

					Expect(page.FindByID("group-" + group_id).Click()).To(Succeed()) // グループを開く
					time.Sleep(time.Second * 1)                                      // 入力画面が出るまで待つ

					Expect(page.FindByID("tag-" + tag_id).Click()).To(Succeed()) // タグをクリック
				})
				It("ツイートが見つかりませんと表示されること", func() {
					Expect(page).To(HaveURL(server_url + "/tweets/" + strings.ToLower(url.PathEscape(tag.Name))))
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
					Expect(page.FindByID("tweet-10")).To(MatchText("tweet-10"))
					Expect(page.FindByID("btn-page-2").Click()).To(Succeed())
					Expect(page.FindByID("tweet-1")).To(MatchText("tweet-1"))
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
		Describe("/:tag", func() {
			var (
				tag    models.Tag
				tag_id string
			)
			BeforeEach(func() {
				db.Find(&tag, "name = ?", "二川二水")
				tag_id = strconv.FormatInt(int64(tag.ID), 10)
			})
			Context("該当するツイートがあるタグで検索した時", func() {
				BeforeEach(func() {
					Expect(page.Navigate(server_url + "/tweets/" + tag.Name)).To(Succeed())
				})
				It("該当ツイートが表示されること", func() {
					var tweet_id = strconv.FormatInt(tweets[len(tweets)-1].ID, 10)
					Expect(page.FindByID("tweet-" + tweet_id)).To(MatchText(tag.Name))
					Expect(page.FindByID("tag-" + tweet_id + "-" + tag_id)).To(MatchText(tag.Name))
					Expect(page.FindByID("btn-page-2").Click()).To(Succeed())

					tweet_id = strconv.FormatInt(tweets[len(tweets)-2].ID, 10)
					Expect(page.FindByID("tweet-" + tweet_id)).To(MatchText(tag.Name))
					Expect(page.FindByID("tag-" + tweet_id + "-" + tag_id)).To(MatchText(tag.Name))
				})
			})
			Context("該当するツイートがないタグで検索した時", func() {
				BeforeEach(func() {
					Expect(page.Navigate(server_url + "/tweets/一柳梨璃")).To(Succeed())
				})
				It("ツイートが見つかりませんと表示されること", func() {
					Expect(page.FindByID("tweet-not-found")).To(HaveText("ツイートが見つかりません"))
				})
			})
		})
	})
})
