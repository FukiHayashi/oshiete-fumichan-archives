package main

import (
	"log"
	"takanome/rareskill"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	rareskill.Takanome()
}
