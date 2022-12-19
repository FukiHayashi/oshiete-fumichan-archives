package rareskill

import (
	"strings"
	"takanome/database"
	"takanome/models"
	"time"

	"gorm.io/gorm"
)

func Register() {
	// DBへ接続
	db := database.DataBaseConnect()
	defer database.DataBaseDisconnect(db)
	// 前回実行日時取得
	registered_at := getLastRegisteredAt(db)
	// 実行日時より新しいtweetを取得
	tweets := getNewTweets(db, registered_at.LastUpdatedAt)
	// 新しいツイートが有れば以下を実行
	if len(tweets) > 0 {
		// Tagを取得
		keywords := getKeywords(db)
		// keywordを含むツイートにTagをつける
		registerTagToTweet(db, tweets, keywords)
	}
	registered_at.LastUpdatedAt = time.Now()
	db.Save(&registered_at)
}

// 前回実行日時取得
func getLastRegisteredAt(db *gorm.DB) models.History {
	var registered_at models.History
	db.Where("name = ?", "registeredAt").First(&registered_at)
	return registered_at
}

// 前回実行日時より新しいtweetを取得
func getNewTweets(db *gorm.DB, registeredAt time.Time) (tweets []models.Tweet) {
	db.Where("created_at > ?", registeredAt).Find(&tweets)
	return tweets
}

// Tagを取得
func getKeywords(db *gorm.DB) (keywords []models.Keyword) {
	db.Preload("Tag").Find(&keywords)
	return keywords
}

// keywordを含むツイートにTagをつける
func registerTagToTweet(db *gorm.DB, tweets []models.Tweet, keywords []models.Keyword) {
	// 更新用スライス
	var update_tweets []models.Tweet
	for _, tweet := range tweets {
		for _, keyword := range keywords {
			// tweetとretweetにキーワードが含まれる場合
			if strings.Contains(tweet.Text, keyword.Name) || strings.Contains(tweet.RetweetText, keyword.Name) {
				// タグが重複する場合は追加しない
				if !tweetContainedTag(tweet.Tags, keyword.Tag) {
					tweet.Tags = append(tweet.Tags, keyword.Tag)
				}
			}
		}
		// 更新するツイートをスライスに追加
		update_tweets = append(update_tweets, tweet)
	}
	// tweetを更新
	if len(update_tweets) > 0 {
		db.Save(&update_tweets)
	}
}

// タグの重複チェック
func tweetContainedTag(tags []models.Tag, newTag models.Tag) bool {
	for _, tag := range tags {
		if tag.ID == newTag.ID {
			return true
		}
	}
	return false
}
