package controllers

import (
	"cineverse/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Admin: Add new movie
func AddMovie(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var movie models.Movie
		if err := c.ShouldBind(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&movie).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin/movies")
	}
}

// User: View all movies
func GetAllMovies(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var movies []models.Movie
		if err := db.Find(&movies).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
			return
		}
		if c.GetHeader("Accept") == "application/json" {
			c.JSON(http.StatusOK, movies)
			return
		}

		c.HTML(http.StatusOK, "movies/movies.html", gin.H{
			"title":      "All Movies",
			"movies":     movies,
			"ActivePage": "movies",
		})
	}
}
