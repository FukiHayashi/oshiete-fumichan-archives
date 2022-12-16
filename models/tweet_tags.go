package models

import (
	"time"

	"gorm.io/gorm"
)

type TweetTags struct {
	TweetID   int64 `gorm:"primaryKey;unique_tag_id"`
	TagID     int   `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
