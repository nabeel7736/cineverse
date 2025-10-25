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
		log.Fatalf("Error Loading .env file: %v",err)
	}
}
