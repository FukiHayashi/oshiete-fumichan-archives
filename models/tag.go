package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique_group_id"`
	GroupID   uint
	Group     Group
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
