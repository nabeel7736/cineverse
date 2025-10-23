package main

import (
	"cineverse/config"
	"cineverse/models"
	"log"
)

func main() {
	config.ConnectDatabase()
	db := config.DB

	// Update user role to admin
	var user models.User
	if err := db.Where("email = ?", "admin@example.com").First(&user).Error; err != nil {
		log.Fatal("User not found")
	}

	user.Role = "admin"
	if err := db.Save(&user).Error; err != nil {
		log.Fatal("Failed to update user role")
	}

	log.Println("User role updated to admin")
}
