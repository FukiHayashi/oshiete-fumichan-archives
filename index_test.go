package main_test

import (
	"takanome/models"
	"takanome/rareskill"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("/", Ordered, func() {
	var (
		page   *agouti.Page
		tweets []models.Tweet
	)
	BeforeAll(func() {
		// ツイートを作成
		for i := 0; i < 10; i++ {
			tweet := tweetFactory.MustCreate().(*models.Tweet)
			tweets = append(tweets, *tweet)
		}
		db.Save(&tweets)
		// タグを登録
		rareskill.Register()
	})
	BeforeEach(func() {
		var err error
		page, err = agouti_driver.NewPage()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		db.Delete(&tweets)
	})
	AfterEach(func() {
		Expect(page.Destroy()).To(Succeed())
	})

	Describe("ページコンテンツ確認", func() {
		BeforeEach(func() {
			Expect(page.Navigate(server_url)).To(Succeed())
		})
		Describe("検索", func() {
			Context("クリックした時", func() {
				BeforeEach(func() {
					Expect(page.FindByID("input-search-keywords-index").Fill("tweet-10")).To(Succeed())
					Expect(page.FindByID("btn-search-index").Click()).To(Succeed())
				})
				It("ツイートを検索できること", func() {
					Expect(page.FindByID("tweet-10")).To(MatchText("tweet-10"))
				})
			})
		})
		Describe("このサイトについて", func() {
			Context("クリックした時", func() {
				BeforeEach(func() {
					Expect(page.FindByID("site-description-btn").Click()).To(Succeed())
				})
				It("このサイトについての説明が表示されること", func() {
					time.Sleep(time.Second * 1) // モーダルが出るまで待つ
					Expect(page.FindByID("site-description-label")).To(MatchText("このサイトについて"))
					Expect(page.FindByID("site-description-close-btn").Click()).To(Succeed())
				})
			})
		})
		Describe("検索対象", func() {
			Context("クリックした時", func() {
				BeforeEach(func() {
					Expect(page.FindByID("search-target-btn").Click()).To(Succeed())
				})
				It("検索対象の説明が表示されること", func() {
					time.Sleep(time.Second * 1) // モーダルが出るまで待つ
					Expect(page.FindByID("search-target-label")).To(MatchText("検索対象"))
					Expect(page.FindByID("search-target-close-btn").Click()).To(Succeed())
				})
			})
		})
	})
	Describe("フッター確認", func() {
		BeforeEach(func() {
			Expect(page.Navigate(server_url)).To(Succeed())
		})
		Describe("利用規約", func() {
			Context("クリックした時", func() {
				BeforeEach(func() {
					Expect(page.FindByID("terms-of-service-link").Click()).To(Succeed())
				})
				It("利用規約の説明が表示されること", func() {
					time.Sleep(time.Second * 1) // モーダルが出るまで待つ
					Expect(page.FindByID("terms-of-service-label")).To(MatchText("利用規約"))
					Expect(page.FindByID("terms-of-service-close-btn").Click()).To(Succeed())
				})
			})
		})
		Describe("プライバシーポリシー", func() {
			Context("クリックした時", func() {
				BeforeEach(func() {
					Expect(page.FindByID("privacy-policy-link").Click()).To(Succeed())
				})
				It("プライバシーポリシーの説明が表示されること", func() {
					time.Sleep(time.Second * 1) // モーダルが出るまで待つ
					Expect(page.FindByID("privacy-policy-label")).To(MatchText("プライバシーポリシー"))
					Expect(page.FindByID("privacy-policy-close-btn").Click()).To(Succeed())
				})
			})
		})
		Describe("ツイートについて", func() {
			Context("クリックした時", func() {
				BeforeEach(func() {
					Expect(page.FindByID("content-of-tweets-link").Click()).To(Succeed())
				})
				It("ツイートについての説明が表示されること", func() {
					time.Sleep(time.Second * 1) // モーダルが出るまで待つ
					Expect(page.FindByID("content-of-tweets-label")).To(MatchText("ツイートについて"))
					Expect(page.FindByID("content-of-tweets-close-btn").Click()).To(Succeed())
				})
			})
		})
	})
	Describe("ヘッダ確認", func() {
		BeforeEach(func() {
			Expect(page.Navigate(server_url)).To(Succeed())
		})
		Describe("ツイート一覧", func() {
			BeforeEach(func() {
				Expect(page.FindByID("tweets-link").Click()).To(Succeed())
			})
			Context("クリックした時", func() {
				It("ツイート一覧ページへ遷移すること", func() {
					Expect(page).To(HaveURL(server_url + "/tweets"))
				})
			})
		})
		Describe("ロゴ", func() {
			BeforeEach(func() {
				Expect(page.FindByID("logo-link").Click()).To(Succeed())
			})
			Context("クリックした時", func() {
				It("トップページへ遷移すること", func() {
					Expect(page).To(HaveURL(server_url + "/"))
				})
			})
		})
	})
})
