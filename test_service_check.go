package main

import (
	"cineverse/config"
	"cineverse/models"
	"cineverse/services"
	"fmt"
)

func main() {
	config.ConnectDatabase()
	var user models.User
	config.DB.Where("email = ?", "admin@example.com").First(&user)
	fmt.Println("PasswordHash:", user.PasswordHash)
	fmt.Println("Role:", user.Role)
	fmt.Println("Check password 'password' with service:", services.CheckPasswordHash("password", user.PasswordHash))
}
