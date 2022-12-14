package models

import (
	"time"

	"gorm.io/gorm"
)

type Keyword struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique_tag_id"`
	TagID     uint
	Tag       Tag
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
