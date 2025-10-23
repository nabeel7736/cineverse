package services

import (
	"cineverse/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RegisterUser(db *gorm.DB, user *models.User) error {
	var existing models.User
	if err := db.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		return errors.New("email already registered")
	}

	hashedpassword, _ := HashPassword(user.PasswordHash)
	user.PasswordHash = hashedpassword

	return db.Create(user).Error
}

func LoginUser(db *gorm.DB, email, password string) (*models.User, error) {
	var user models.User
	if err := db.Where("email =?", email).First(&user).Error; err != nil {
		return nil, errors.New("Inavlid email or Password")
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("Invalid email or password")
	}
	return &user, nil
}
