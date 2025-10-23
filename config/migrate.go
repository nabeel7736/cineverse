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
		&models.Movie{},
		&models.Booking{},
	)
	if err != nil {
		fmt.Println("Migration Failed", err)
		return
	}
	fmt.Println("Migrated Successfully !!")

}
