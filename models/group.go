package models

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"unique_category_id"`
	CategoryID uint
	Category   Category
	Tags       []Tag
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
