package database

import (
	"database/sql"
	"log"

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
