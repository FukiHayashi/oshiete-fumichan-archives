package database

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"log"
	"takanome/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"os"
)

func DataBaseConnect() *gorm.DB {
	// Connect to Database
	sqlDB, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	return db
}

func DataBaseDisconnect(db *gorm.DB) {
	dbc, _ := db.DB()
	dbc.Close()
}

func DataBaseInit() {
	db := DataBaseConnect()
	defer DataBaseDisconnect(db)

	// オートマイグレーション
	db.AutoMigrate(&models.Tweet{}, &models.Category{}, &models.Tag{}, &models.Keyword{}, &models.History{})

	// 実行履歴がなければDBへ初期データを投入する
	var registered_at models.History
	if err := db.Where("name = ?", "registeredAt").First(&registered_at).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("実行履歴が見つからないためDBを初期化します。")
		// categoryの初期化
		categories := categoriesInit()
		db.Save(&categories)
		// groupsの初期化
		groups := groupsInit(categories)
		db.Save(&groups)
		// tagsの初期化
		tags := tagsInit(groups)
		db.Save(tags)
		// keywordsの初期化
		keywords := keywordsInit(tags)
		db.Save(keywords)
		// 実行履歴
		registered_at.Name = "registeredAt"
		registered_at.LastUpdatedAt = time.Now().AddDate(-20, 0, 0)
		db.Save(&registered_at)
	}
}

func categoriesInit() (categories []models.Category) {
	rows := loadCSV(os.Getenv("CATEGORIES_CSV"))
	for _, row := range rows {
		if row[0] != "" {
			category := models.Category{
				Name: row[0],
			}
			categories = append(categories, category)
		}
	}
	return categories
}

func groupsInit(categories []models.Category) (groups []models.Group) {
	rows := loadCSV(os.Getenv("GROUPS_CSV"))
	for _, row := range rows {
		if row[0] != "" {
			for _, category := range categories {
				if category.Name == row[0] {
					group := models.Group{
						Name:       row[1],
						CategoryID: category.ID,
						Category:   category,
					}
					groups = append(groups, group)
				}
			}
		}
	}
	return groups
}

func tagsInit(groups []models.Group) (tags []models.Tag) {
	rows := loadCSV(os.Getenv("TAGS_CSV"))
	for _, row := range rows {
		if row[0] != "" {
			for _, group := range groups {
				if group.Name == row[0] {
					tag := models.Tag{
						Name:    row[1],
						GroupID: group.ID,
						Group:   group,
					}
					tags = append(tags, tag)
				}
			}
		}
	}
	return tags
}

func keywordsInit(tags []models.Tag) (keywords []models.Keyword) {
	rows := loadCSV(os.Getenv("KEYWORDS_CSV"))
	for _, row := range rows {
		if row[0] != "" {
			for _, tag := range tags {
				if tag.Name == row[0] {
					for i := 1; i < len(row); i++ {
						if row[i] != "" {
							keyword := models.Keyword{
								Name:  row[i],
								TagID: tag.ID,
								Tag:   tag,
							}
							keywords = append(keywords, keyword)
						}
					}
				}
			}
		}
	}
	return keywords
}

// csv
func loadCSV(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return rows
}
