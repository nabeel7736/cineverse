package controllers

import (
	"cineverse/models"
	"cineverse/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ShowRegisterpage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{"title": "Register"})
}

func ShowLoginPage(c *gin.Context) {
	message := ""
	if c.Query("registered") == "1" {
		message = "Registration successful. Please log in."
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":      "Login",
		"message":    message,
		"ActivePage": "login",
	})
}

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fullname := strings.TrimSpace(ctx.PostForm("full_name"))
		email := strings.ToLower(strings.TrimSpace(ctx.PostForm("email")))
		password := strings.TrimSpace(ctx.PostForm("password"))

		if fullname == "" || email == "" || password == "" {
			ctx.HTML(http.StatusBadRequest, "register.html", gin.H{
				"title":   "Register",
				"error":   "All fields are required.",
				"message": "",
			})
			return
		}
		// Check duplicate email
		var existing models.User
		if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
			ctx.HTML(http.StatusConflict, "register.html", gin.H{
				"error":      "Email already registered.",
				"ActivePage": "register",
			})
			return
		}
		hashed, err := utils.HashPassword(password)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "register.html", gin.H{
				"error":      "Failed to hash password.",
				"ActivePage": "register",
			})
			return
		}
		user := models.User{
			FullName:     fullname,
			Email:        email,
			PasswordHash: hashed,
			Role:         "user",
			IsVerified:   false,
			IsBLocked:    false,
		}

		if err := db.Create(&user).Error; err != nil {
			ctx.HTML(http.StatusInternalServerError, "register.html", gin.H{
				"error":      "Failed to create user.",
				"ActivePage": "register",
			})
			return
		}

		ctx.Redirect(http.StatusSeeOther, "/login?registered=1")

	}
}

// func LoginHandler(db *gorm.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		email := strings.TrimSpace(ctx.PostForm("email"))
// 		password := strings.TrimSpace(ctx.PostForm("password"))

// 		if email == "" || password == "" {
// 			ctx.HTML(http.StatusBadRequest, "login.html", gin.H{
// 				"title": "Login Page",
// 				"error": "Email and Password cannot be empty",
// 			})
// 			return
// 		}

// 		if !utils.IsValidEmail(email) {
// 			ctx.HTML(http.StatusBadRequest, "login.html", gin.H{
// 				"title": "Login Page",
// 				"error": "Invalid email format",
// 			})
// 			return
// 		}
// 		var user models.User
// 		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
// 			ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{
// 				"error":      "Email not registered.",
// 				"ActivePage": "login",
// 			})
// 			return
// 		}

// 		if !utils.CheckPasswordHash(user.PasswordHash, password) {
// 			ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{
// 				"error":      "Invalid password.",
// 				"ActivePage": "login",
// 			})
// 			return
// 		}

// 		if user.IsBLocked {
// 			ctx.HTML(http.StatusForbidden, "login.html", gin.H{
// 				"error":      "Account blocked by admin.",
// 				"ActivePage": "login",
// 			})
// 			return
// 		}

// 		// ================== Create Access & Refresh Tokens ==================
// 		accessToken, err := utils.CreateAccessToken(user.ID)
// 		if err != nil {
// 			ctx.HTML(http.StatusInternalServerError, "login.html", gin.H{
// 				"error":      "Failed to create access token.",
// 				"ActivePage": "login",
// 			})
// 			return
// 		}

// 		refreshToken, err := utils.CreateRefreshToken(user.ID)
// 		if err != nil {
// 			ctx.HTML(http.StatusInternalServerError, "login.html", gin.H{
// 				"error":      "Failed to create refresh token.",
// 				"ActivePage": "login",
// 			})
// 			return
// 		}

// 		user.RefreshToken = refreshToken
// 		if err := db.Save(&user).Error; err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
// 			return
// 		}

// 		ctx.SetCookie("access_token", accessToken, 3600, "/", "", false, true) // 1 hour
// 		// ctx.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true) // 7 days
// 		ctx.SetCookie("user_name", user.FullName, 3600*24*7, "/", "", false, false)

// 		// Redirect based on role
// 		if user.Role == "admin" {
// 			ctx.Redirect(http.StatusSeeOther, "/admin/dashboard")
// 			return
// 		}

// 		// ctx.Redirect(http.StatusSeeOther, "/movies")
// 		ctx.JSON(http.StatusOK, gin.H{"message": "login successfull !!!"})
// 	}
// }

// func Logout(c *gin.Context) {
// 	c.SetCookie("access_token", "", -1, "/", "", false, true)
// 	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
// 	c.SetCookie("user_name", "", -1, "/", "", false, false)
// 	// c.Redirect(http.StatusSeeOther, "/login")

// 	c.JSON(http.StatusOK, gin.H{"message": "Logout sucessfull"})
// }

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := strings.TrimSpace(ctx.PostForm("email"))
		password := strings.TrimSpace(ctx.PostForm("password"))

		var user models.User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
			return
		}

		if !utils.CheckPasswordHash(user.PasswordHash, password) {
			ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid password"})
			return
		}

		// Generate tokens
		accessToken, _ := utils.CreateAccessToken(user.ID)
		refreshToken, _ := utils.CreateRefreshToken(user.ID)

		// Save refresh token in DB
		// tokenEntry := models.RefreshToken{
		// 	UserID:    user.ID,
		// 	Token:     refreshToken,
		// 	ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		// }
		// db.Create(&tokenEntry)

		db.Create(&models.RefreshToken{
			UserID:    user.ID,
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		})

		// Cookies
		ctx.SetCookie("access_token", accessToken, 3600, "/", "", false, true)
		ctx.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)
		ctx.SetCookie("user_name", user.FullName, 3600*24*7, "/", "", false, false)

		if user.Role == "admin" {
			ctx.Redirect(http.StatusSeeOther, "/admin/dashboard")
			return
		}

		ctx.Redirect(http.StatusSeeOther, "/movies")

		ctx.JSON(http.StatusOK, gin.H{"message": "Login successful", "access_token": accessToken})
	}
}

// func LogoutHandler(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Get user_id from middleware (AuthMiddleware must set it)
// 		userID, exists := c.Get("user_id")
// 		if exists {
// 			// Clear the refresh token in DB to invalidate it
// 			db.Model(&models.User{}).Where("id = ?", userID).Update("refresh_token", "")
// 		}

// 		// Delete cookies (set expiry to -1)
// 		c.SetCookie("access_token", "", -1, "/", "", false, true)
// 		c.SetCookie("refresh_token", "", -1, "/", "", false, true)
// 		c.SetCookie("user_name", "", -1, "/", "", false, false)

// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "Logout successful. Tokens cleared and session ended.",
// 		})
// 	}
// }

// func RefreshTokenHandler(db *gorm.DB) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		refreshToken, err := ctx.Cookie("refresh_token")
// 		if err != nil || refreshToken == "" {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
// 			return
// 		}

// 		claims, err := utils.ValidateRefreshToken(refreshToken)
// 		if err != nil {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
// 			return
// 		}

// 		var user models.User
// 		if err := db.First(&user, claims["user_id"]).Error; err != nil {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
// 			return
// 		}

// 		if user.RefreshToken != refreshToken {
// 			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token revoked"})
// 			return
// 		}

// 		// create new access token
// 		newAccessToken, err := utils.CreateAccessToken(user.ID)
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
// 			return
// 		}

// 		ctx.SetCookie("access_token", newAccessToken, 3600, "/", "", false, true)
// 		ctx.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
// 	}
// }

func RefreshTokenHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		refreshToken, err := ctx.Cookie("refresh_token")
		if err != nil || refreshToken == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
			return
		}

		claims, err := utils.ValidateRefreshToken(refreshToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		userID := uint(claims["user_id"].(float64))

		// var dbToken models.RefreshToken
		// if err := db.Where("user_id = ? AND token = ?", userID, refreshToken).First(&dbToken).Error; err != nil {
		// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		// 	return
		// }

		var dbToken models.RefreshToken
		if err := db.Where("user_id = ? AND token = ?", userID, refreshToken).First(&dbToken).Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
			return
		}

		if time.Now().After(dbToken.ExpiresAt) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
			db.Delete(&dbToken)
			return
		}

		newAccessToken, _ := utils.CreateAccessToken(userID)
		ctx.SetCookie("access_token", newAccessToken, 3600, "/", "", false, true)
		ctx.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
	}
}

func LogoutHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("user_id")
		if exists {
			db.Where("user_id = ?", userID).Delete(&models.RefreshToken{})
		}

		ctx.SetCookie("access_token", "", -1, "/", "", false, true)
		ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)
		ctx.SetCookie("user_name", "", -1, "/", "", false, false)

		ctx.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	}
}
