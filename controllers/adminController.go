package controllers

import (
	"cineverse/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminDashboard(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var totalMovies int64
		var totalBookings int64
		var totalUsers int64

		db.Model(&models.Movie{}).Count(&totalMovies)
		db.Model(&models.Booking{}).Count(&totalBookings)
		db.Model(&models.User{}).Count(&totalUsers)

		c.HTML(http.StatusOK, "base.html", gin.H{
			"title":         "Admin Dashboard",
			"TotalMovies":   totalMovies,
			"TotalBookings": totalBookings,
			"TotalUsers":    totalUsers,
			"ActivePage":    "admin",
		})
	}
}
