package models

import (
	"time"

	"gorm.io/gorm"
)

type History struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"unique"`
	LastUpdatedAt time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
