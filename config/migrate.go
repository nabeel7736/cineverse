package config

import (
	"cineverse/models"
	"fmt"
)

func MigrateAll() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.OTP{},
		&models.RefreshToken{},
	)
	if err != nil {
		fmt.Println("Migration Failed", err)
		return
	}
	fmt.Println("Migrated Successfully !!")

}
