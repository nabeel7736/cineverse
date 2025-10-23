package routes

import (
	"cineverse/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MovieRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// User accessible
	// rg.GET("/movies", middlewares.AuthMiddleware(), controllers.GetAllMovies(db))

	// Admin accessible (later protect with admin middleware)
	rg.POST("/admin/movies/add", controllers.AddMovie(db))
}
