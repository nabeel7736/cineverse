package main

import (
	"cineverse/config"
	"cineverse/models"
	"cineverse/utils"
	"fmt"
)

func main() {
	config.ConnectDatabase()
	var user models.User
	config.DB.Where("email = ?", "admin@example.com").First(&user)
	fmt.Println("Old PasswordHash:", user.PasswordHash)
	newHash, _ := utils.HashPassword("password")
	user.PasswordHash = newHash
	config.DB.Save(&user)
	fmt.Println("New PasswordHash:", user.PasswordHash)
}
