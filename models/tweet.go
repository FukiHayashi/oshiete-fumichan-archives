package models

import (
	"time"

	"gorm.io/gorm"
)

type Tweet struct {
	// ツイートの情報
	ID         int64     `gorm:"primaryKey;unique;not null" json:"id"`
	Text       string    `gorm:"not null" json:"text"`
	TweetedAt  time.Time `gorm:"not null" json:"tweeted_at"`
	Url        string    `gorm:"not null" json:"url"`
	RawData    string    `gorm:"not null" json:"raw_data"`
	ScreenName string    `gorm:"not null" json:"screen_name"`
	// 質問ツイートの情報
	RetweetScreenName string `json:"retweet_screen_name"`
	RetweetText       string `json:"retweet_text"`
	RetweetUrl        string `json:"retweet_url"`
	// タグ
	Tags []Tag `gorm:"many2many:tweet_tags"`
	// DB情報
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
