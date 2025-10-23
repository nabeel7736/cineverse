package config

import (
	"log"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Loadenv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}
}
