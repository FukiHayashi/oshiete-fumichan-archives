package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	CategoryID uint
	Category   Category
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
