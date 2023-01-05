package controllers

import (
	"takanome/models"

	"gorm.io/gorm"
)

// 全てのカテゴリー、グループ、タグを取得
func GetAllCategories(db *gorm.DB) []models.Category {
	var categories []models.Category

	db.Preload("Groups.Tags").Find(&categories)

	return categories
}
